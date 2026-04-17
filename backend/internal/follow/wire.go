//go:build wireinject
// +build wireinject

package follow

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/follow/application"
	"github.com/blog/blog-community/internal/follow/delivery"
	"github.com/blog/blog-community/internal/follow/infrastructure"
)

// InitializeHandler injects follow dependencies.
func InitializeHandler(client *ent.Client) *delivery.FollowHandler {
	wire.Build(
		infrastructure.NewEntFollowRepo,
		application.NewFollowUseCase,
		delivery.NewFollowHandler,
	)
	return &delivery.FollowHandler{}
}
