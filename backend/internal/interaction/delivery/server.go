package delivery

import "github.com/gin-gonic/gin"

// InteractionServer groups all HTTP handlers for the interaction module.
type InteractionServer struct {
	handler *InteractionHandler
}

// NewInteractionServer creates an interaction delivery server.
func NewInteractionServer(handler *InteractionHandler) *InteractionServer {
	return &InteractionServer{handler: handler}
}

// Register mounts interaction routes on the given router group.
func (s *InteractionServer) Register(r *gin.RouterGroup, auth gin.HandlerFunc) {
	r.GET("/likes/:targetType/:targetId", s.handler.GetLikeStatus)
	r.GET("/collects/:targetType/:targetId", s.handler.GetCollectStatus)

	authGroup := r.Group("/")
	authGroup.Use(auth)
	authGroup.POST("/likes/:targetType/:targetId", s.handler.ToggleLike)
	authGroup.POST("/collects/:targetType/:targetId", s.handler.ToggleCollect)
}
