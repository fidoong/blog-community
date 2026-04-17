//go:build wireinject
// +build wireinject

package feed

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/feed/application"
	"github.com/blog/blog-community/internal/feed/delivery"
	"github.com/blog/blog-community/internal/feed/domain"
	postapp "github.com/blog/blog-community/internal/post/application"
	postinfra "github.com/blog/blog-community/internal/post/infrastructure"
)

// InitializeHandler injects feed dependencies.
func InitializeHandler(client *ent.Client) *delivery.FeedHandler {
	wire.Build(
		postinfra.NewEntPostRepo,
		postapp.NewPostUseCase,
		wire.Bind(new(domain.PostLister), new(postapp.UseCase)),
		application.NewFeedUseCase,
		delivery.NewFeedHandler,
	)
	return &delivery.FeedHandler{}
}
