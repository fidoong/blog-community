//go:build wireinject
// +build wireinject

package post

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/post/application"
	"github.com/blog/blog-community/internal/post/delivery"
	"github.com/blog/blog-community/internal/post/infrastructure"
	"github.com/blog/blog-community/internal/ent"
	"github.com/redis/go-redis/v9"
)

// InitializeHandler wires all dependencies for the post HTTP handler.
func InitializeHandler(client *ent.Client, redisClient *redis.Client) *delivery.PostHandler {
	wire.Build(
		infrastructure.NewEntPostRepo,
		infrastructure.NewSearchTrendRepo,
		application.NewPostUseCase,
		delivery.NewPostHandler,
	)
	return &delivery.PostHandler{}
}
