package delivery

import "github.com/gin-gonic/gin"

// Server groups all HTTP handlers for the user service.
type Server struct {
	User  *UserHandler
	OAuth *OAuthHandler
}

// NewServer creates a delivery server from handlers.
func NewServer(user *UserHandler, oauth *OAuthHandler) *Server {
	return &Server{User: user, OAuth: oauth}
}

// Register mounts routes on the given router group.
func (s *Server) Register(r *gin.RouterGroup) {
	r.POST("/auth/register", s.User.Register)
	r.POST("/auth/login", s.User.Login)
	r.GET("/users/:id", s.User.GetProfile)

	r.GET("/auth/oauth/github", s.OAuth.GetGitHubAuthURL)
	r.GET("/auth/oauth/google", s.OAuth.GetGoogleAuthURL)
	r.GET("/auth/oauth/github/callback", s.OAuth.GitHubCallback)
	r.GET("/auth/oauth/google/callback", s.OAuth.GoogleCallback)
}
