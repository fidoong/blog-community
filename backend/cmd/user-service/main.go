package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/blog/blog-community/configs"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/cache"
	"github.com/blog/blog-community/pkg/logger"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/internal/comment/application"
	commentDelivery "github.com/blog/blog-community/internal/comment/delivery"
	commentInfra "github.com/blog/blog-community/internal/comment/infrastructure"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/migrate"
	interactionApp "github.com/blog/blog-community/internal/interaction/application"
	interactionDelivery "github.com/blog/blog-community/internal/interaction/delivery"
	interactionInfra "github.com/blog/blog-community/internal/interaction/infrastructure"
	postApplication "github.com/blog/blog-community/internal/post/application"
	postDelivery "github.com/blog/blog-community/internal/post/delivery"
	postInfra "github.com/blog/blog-community/internal/post/infrastructure"
	"github.com/blog/blog-community/internal/user"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Init(cfg.AppEnv)
	log := logger.L()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed opening connection to postgres", zap.Error(err))
	}
	defer client.Close()

	// Run auto migration
	ctx := context.Background()
	if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		log.Fatal("failed creating schema resources", zap.Error(err))
	}

	// Connect to Redis
	redisClient := cache.NewRedisClient(cfg.RedisAddr)
	defer redisClient.Close()
	tokenStore := auth.NewRedisTokenStore(redisClient)

	// Wire up handler
	server := user.InitializeServer(cfg, client, tokenStore)

	// Setup router
	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// Health checks
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	authMiddleware := middleware.AuthRequired(cfg.JWTSecret)
	server.Register(api, authMiddleware)

	// Post routes
	postRepo := postInfra.NewEntPostRepo(client)
	postUseCase := postApplication.NewPostUseCase(postRepo)
	postHandler := postDelivery.NewPostHandler(postUseCase)
	api.GET("/posts", postHandler.List)
	api.GET("/posts/:id", postHandler.GetByID)
	api.POST("/posts", authMiddleware, postHandler.Create)
	api.PUT("/posts/:id", authMiddleware, postHandler.Update)
	api.DELETE("/posts/:id", authMiddleware, postHandler.Delete)
	api.POST("/posts/:id/publish", authMiddleware, postHandler.Publish)

	// Comment routes
	commentRepo := commentInfra.NewEntCommentRepo(client)
	commentUseCase := application.NewCommentUseCase(commentRepo)
	commentHandler := commentDelivery.NewCommentHandler(commentUseCase)
	api.GET("/posts/:id/comments", commentHandler.List)
	api.POST("/posts/:id/comments", authMiddleware, commentHandler.Create)
	api.DELETE("/comments/:id", authMiddleware, commentHandler.Delete)

	// Interaction routes
	interactionRepo := interactionInfra.NewEntInteractionRepo(client)
	interactionCounter := interactionInfra.NewRedisCounter(redisClient)
	interactionUseCase := interactionApp.NewInteractionUseCase(interactionRepo, interactionCounter)
	interactionHandler := interactionDelivery.NewInteractionHandler(interactionUseCase)
	api.POST("/likes/:targetType/:targetId", authMiddleware, interactionHandler.ToggleLike)
	api.GET("/likes/:targetType/:targetId", interactionHandler.GetLikeStatus)
	api.POST("/collects/:targetType/:targetId", authMiddleware, interactionHandler.ToggleCollect)
	api.GET("/collects/:targetType/:targetId", interactionHandler.GetCollectStatus)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Info("starting user service", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}
	log.Info("server exited")
}
