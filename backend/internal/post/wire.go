//go:build wireinject
// +build wireinject

package post

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/internal/post/application"
	"github.com/blog/blog-community/internal/post/delivery"
	"github.com/blog/blog-community/internal/post/infrastructure"
	"github.com/blog/blog-community/internal/ent"
)

// InitializeHandler wires all dependencies for the post HTTP handler.
func InitializeHandler(client *ent.Client) *delivery.PostHandler {
	wire.Build(
		infrastructure.NewEntPostRepo,
		application.NewPostUseCase,
		delivery.NewPostHandler,
	)
	return &delivery.PostHandler{}
}
