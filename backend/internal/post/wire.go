//go:build wireinject
// +build wireinject

package post

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/post/application"
	"github.com/blog/blog-community/internal/post/delivery"
	"github.com/blog/blog-community/internal/post/infrastructure"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/pkg/search"
	"github.com/redis/go-redis/v9"
)

// InitializeHandler wires all dependencies for the post HTTP handler.
func InitializeHandler(client *ent.Client, redisClient *redis.Client, esClient *search.Client) *delivery.PostHandler {
	wire.Build(
		infrastructure.NewEntPostRepo,
		infrastructure.NewSearchTrendRepo,
		wire.Bind(new(application.SearchTrendRecorder), new(*infrastructure.SearchTrendRepo)),
		infrastructure.NewPostIndexer,
		wire.Bind(new(application.PostIndexer), new(*infrastructure.PostIndexer)),
		application.NewPostUseCase,
		delivery.NewPostHandler,
	)
	return &delivery.PostHandler{}
}
