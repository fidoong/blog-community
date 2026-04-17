package delivery

import "github.com/gin-gonic/gin"

// Server is the feed module route registrar.
type Server struct {
	Feed *FeedHandler
}

// NewFeedServer creates a new feed server.
func NewFeedServer(feed *FeedHandler) *Server {
	return &Server{Feed: feed}
}

// Register registers feed routes.
func (s *Server) Register(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	r.GET("/feed", s.Feed.GetFeed)
}
