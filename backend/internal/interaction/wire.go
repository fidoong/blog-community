//go:build wireinject
// +build wireinject

package interaction

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/interaction/application"
	"github.com/blog/blog-community/internal/interaction/delivery"
	"github.com/blog/blog-community/internal/interaction/infrastructure"
	notificationDomain "github.com/blog/blog-community/internal/notification/domain"
	"github.com/redis/go-redis/v9"
)

// InitializeHandler wires all dependencies for the interaction HTTP handler.
func InitializeHandler(client *ent.Client, redisClient *redis.Client, notifier notificationDomain.Notifier) *delivery.InteractionHandler {
	wire.Build(
		infrastructure.NewEntInteractionRepo,
		infrastructure.NewRedisCounter,
		application.NewInteractionUseCase,
		delivery.NewInteractionHandler,
	)
	return &delivery.InteractionHandler{}
}
