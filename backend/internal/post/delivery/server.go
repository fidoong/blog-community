package delivery

import "github.com/gin-gonic/gin"

// PostServer groups all HTTP handlers for the post module.
type PostServer struct {
	handler *PostHandler
}

// NewPostServer creates a post delivery server.
func NewPostServer(handler *PostHandler) *PostServer {
	return &PostServer{handler: handler}
}

// Register mounts post routes on the given router group.
func (s *PostServer) Register(r *gin.RouterGroup, auth gin.HandlerFunc) {
	r.GET("/posts", s.handler.List)
	r.GET("/posts/:id", s.handler.GetByID)
	r.GET("/posts/:id/related", s.handler.GetRelated)
	r.GET("/search/hot", s.handler.HotKeywords)
	r.GET("/search", s.handler.Search)

	authGroup := r.Group("/")
	authGroup.Use(auth)
	authGroup.POST("/posts", s.handler.Create)
	authGroup.PUT("/posts/:id", s.handler.Update)
	authGroup.DELETE("/posts/:id", s.handler.Delete)
	authGroup.POST("/posts/:id/publish", s.handler.Publish)
}
