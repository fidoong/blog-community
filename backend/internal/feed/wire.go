//go:build wireinject
// +build wireinject

package feed

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/feed/application"
	"github.com/blog/blog-community/internal/feed/delivery"
	"github.com/blog/blog-community/internal/feed/domain"
	followdomain "github.com/blog/blog-community/internal/follow/domain"
	followinfra "github.com/blog/blog-community/internal/follow/infrastructure"
	postapp "github.com/blog/blog-community/internal/post/application"
	postinfra "github.com/blog/blog-community/internal/post/infrastructure"
	"github.com/redis/go-redis/v9"
)

func nilPostIndexer() postapp.PostIndexer {
	return nil
}

// InitializeHandler injects feed dependencies.
func InitializeHandler(client *ent.Client, redisClient *redis.Client) *delivery.FeedHandler {
	wire.Build(
		postinfra.NewEntPostRepo,
		postinfra.NewSearchTrendRepo,
		wire.Bind(new(postapp.SearchTrendRecorder), new(*postinfra.SearchTrendRepo)),
		nilPostIndexer,
		postapp.NewPostUseCase,
		wire.Bind(new(domain.PostLister), new(postapp.UseCase)),
		followinfra.NewEntFollowRepo,
		wire.Bind(new(domain.FollowLister), new(followdomain.Repository)),
		application.NewFeedUseCase,
		delivery.NewFeedHandler,
	)
	return &delivery.FeedHandler{}
}
