//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"
	"github.com/blog/blog-community/configs"
	"github.com/blog/blog-community/internal/user/application"
	"github.com/blog/blog-community/internal/user/delivery"
	"github.com/blog/blog-community/internal/user/ent"
	"github.com/blog/blog-community/internal/user/infrastructure"
)

// InitializeServer wires all dependencies for the user HTTP server.
func InitializeServer(cfg *configs.Config, client *ent.Client) *delivery.Server {
	wire.Build(
		infrastructure.NewEntUserRepo,
		infrastructure.NewEntTransactor,
		application.NewUserUseCase,
		delivery.NewUserHandlerFromConfig,
		delivery.NewOAuthHandlerFromConfig,
		delivery.NewServer,
	)
	return &delivery.Server{}
}
