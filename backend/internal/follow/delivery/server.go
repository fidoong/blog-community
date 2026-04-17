package delivery

import "github.com/gin-gonic/gin"

// Server is the follow module route registrar.
type Server struct {
	Follow *FollowHandler
}

// NewFollowServer creates a new follow server.
func NewFollowServer(follow *FollowHandler) *Server {
	return &Server{Follow: follow}
}

// Register registers follow routes.
func (s *Server) Register(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Public
	r.GET("/users/:id/followers", s.Follow.ListFollowers)
	r.GET("/users/:id/following", s.Follow.ListFollowing)
	r.GET("/users/:id/follow-stats", s.Follow.GetFollowStats)

	// Auth required
	authGroup := r.Group("/")
	authGroup.Use(authMiddleware)
	authGroup.POST("/users/:id/follow", s.Follow.Follow)
	authGroup.DELETE("/users/:id/follow", s.Follow.Unfollow)
}
