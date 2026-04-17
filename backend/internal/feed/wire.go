//go:build wireinject
// +build wireinject

package feed

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/feed/application"
	"github.com/blog/blog-community/internal/feed/delivery"
	"github.com/blog/blog-community/internal/feed/domain"
	followinfra "github.com/blog/blog-community/internal/follow/infrastructure"
	postapp "github.com/blog/blog-community/internal/post/application"
	postinfra "github.com/blog/blog-community/internal/post/infrastructure"
	"github.com/redis/go-redis/v9"
)

// InitializeHandler injects feed dependencies.
func InitializeHandler(client *ent.Client, redisClient *redis.Client) *delivery.FeedHandler {
	wire.Build(
		postinfra.NewEntPostRepo,
		postinfra.NewSearchTrendRepo,
		postapp.NewPostUseCase,
		wire.Bind(new(domain.PostLister), new(postapp.UseCase)),
		followinfra.NewEntFollowRepo,
		wire.Bind(new(domain.FollowLister), new(*followinfra.EntFollowRepo)),
		application.NewFeedUseCase,
		delivery.NewFeedHandler,
	)
	return &delivery.FeedHandler{}
}
