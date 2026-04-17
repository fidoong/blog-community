//go:build wireinject
// +build wireinject

package comment

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/comment/application"
	"github.com/blog/blog-community/internal/comment/delivery"
	"github.com/blog/blog-community/internal/comment/infrastructure"
	"github.com/blog/blog-community/internal/ent"
	notificationDomain "github.com/blog/blog-community/internal/notification/domain"
)

// InitializeHandler wires all dependencies for the comment HTTP handler.
func InitializeHandler(client *ent.Client, notifier notificationDomain.Notifier) *delivery.CommentHandler {
	wire.Build(
		infrastructure.NewEntCommentRepo,
		application.NewCommentUseCase,
		delivery.NewCommentHandler,
	)
	return &delivery.CommentHandler{}
}
