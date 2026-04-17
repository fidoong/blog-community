//go:build wireinject
// +build wireinject

package notification

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/notification/application"
	"github.com/blog/blog-community/internal/notification/delivery"
	"github.com/blog/blog-community/internal/notification/infrastructure"
)

// InitializeHandler injects notification dependencies.
func InitializeHandler(client *ent.Client) *delivery.NotificationHandler {
	wire.Build(
		infrastructure.NewEntNotificationRepo,
		application.NewNotificationUseCase,
		delivery.NewNotificationHandler,
	)
	return &delivery.NotificationHandler{}
}

// InitializeUseCase injects notification usecase only.
func InitializeUseCase(client *ent.Client) application.UseCase {
	wire.Build(
		infrastructure.NewEntNotificationRepo,
		application.NewNotificationUseCase,
	)
	return nil
}
