package delivery

import "github.com/gin-gonic/gin"

// CommentServer groups all HTTP handlers for the comment module.
type CommentServer struct {
	handler *CommentHandler
}

// NewCommentServer creates a comment delivery server.
func NewCommentServer(handler *CommentHandler) *CommentServer {
	return &CommentServer{handler: handler}
}

// Register mounts comment routes on the given router group.
func (s *CommentServer) Register(r *gin.RouterGroup, auth gin.HandlerFunc) {
	r.GET("/posts/:postId/comments", s.handler.List)

	authGroup := r.Group("/")
	authGroup.Use(auth)
	authGroup.POST("/posts/:postId/comments", s.handler.Create)
	authGroup.DELETE("/comments/:id", s.handler.Delete)
}
