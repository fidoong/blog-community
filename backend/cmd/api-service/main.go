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

	"github.com/blog/blog-community/configs"
	"github.com/blog/blog-community/internal/comment"
	commentDelivery "github.com/blog/blog-community/internal/comment/delivery"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/migrate"
	"github.com/blog/blog-community/internal/feed"
	feedDelivery "github.com/blog/blog-community/internal/feed/delivery"
	"github.com/blog/blog-community/internal/follow"
	followDelivery "github.com/blog/blog-community/internal/follow/delivery"
	"github.com/blog/blog-community/internal/interaction"
	interactionDelivery "github.com/blog/blog-community/internal/interaction/delivery"
	"github.com/blog/blog-community/internal/notification"
	notificationDelivery "github.com/blog/blog-community/internal/notification/delivery"
	"github.com/blog/blog-community/internal/post"
	postDelivery "github.com/blog/blog-community/internal/post/delivery"
	"github.com/blog/blog-community/internal/user"
	"github.com/blog/blog-community/pkg/auth"
	"github.com/blog/blog-community/pkg/cache"
	"github.com/blog/blog-community/pkg/logger"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/search"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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
	redisClient := cache.NewRedisClientWithPassword(cfg.RedisAddr, cfg.RedisPassword)
	defer redisClient.Close()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("failed connecting to redis", zap.Error(err))
	}
	tokenStore := auth.NewRedisTokenStore(redisClient)

	// Connect to Elasticsearch (optional, log warning if unavailable)
	var esClient *search.Client
	if cfg.ElasticsearchAddr != "" {
		esClient, err = search.NewClient([]string{cfg.ElasticsearchAddr})
		if err != nil {
			log.Warn("failed to create elasticsearch client, search will fallback to DB", zap.Error(err))
		} else if err := esClient.Ping(ctx); err != nil {
			log.Warn("elasticsearch unavailable, search will fallback to DB", zap.Error(err))
			esClient = nil
		} else {
			log.Info("elasticsearch connected", zap.String("addr", cfg.ElasticsearchAddr))
		}
	}

	// Initialize notification usecase first (used as notifier by other modules)
	notificationUC := notification.InitializeUseCase(client)

	// Wire up handlers — each module owns its own route registration
	userServer := user.InitializeServer(cfg, client, tokenStore)
	postHandler := post.InitializeHandler(client, redisClient, esClient)
	commentHandler := comment.InitializeHandler(client, notificationUC)
	interactionHandler := interaction.InitializeHandler(client, redisClient, notificationUC)
	followHandler := follow.InitializeHandler(client, notificationUC)

	// Reindex published posts to ES on startup
	if esClient != nil {
		log.Info("reindexing published posts to elasticsearch...")
		if err := postHandler.ReindexSearch(ctx); err != nil {
			log.Warn("failed to reindex posts", zap.Error(err))
		} else {
			log.Info("elasticsearch reindex completed")
		}
	}

	postServer := postDelivery.NewPostServer(postHandler)
	feedServer := feedDelivery.NewFeedServer(feed.InitializeHandler(client, redisClient))
	followServer := followDelivery.NewFollowServer(followHandler)
	commentServer := commentDelivery.NewCommentServer(commentHandler)
	interactionServer := interactionDelivery.NewInteractionServer(interactionHandler)
	notificationServer := notificationDelivery.NewNotificationServer(notification.InitializeHandler(client))

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

	// API routes — each module registers itself
	api := r.Group("/api/v1")
	authMiddleware := middleware.AuthRequired(cfg.JWTSecret)
	userServer.Register(api, authMiddleware)
	postServer.Register(api, authMiddleware)
	feedServer.Register(api, authMiddleware)
	followServer.Register(api, authMiddleware)
	commentServer.Register(api, authMiddleware)
	interactionServer.Register(api, authMiddleware)
	notificationServer.Register(api, authMiddleware)

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
