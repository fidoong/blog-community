package delivery

import "github.com/gin-gonic/gin"

// Server is the notification module route registrar.
type Server struct {
	Notification *NotificationHandler
}

// NewNotificationServer creates a new notification server.
func NewNotificationServer(notification *NotificationHandler) *Server {
	return &Server{Notification: notification}
}

// Register registers notification routes.
func (s *Server) Register(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	authGroup := r.Group("/")
	authGroup.Use(authMiddleware)
	authGroup.GET("/notifications", s.Notification.List)
	authGroup.GET("/notifications/unread-count", s.Notification.CountUnread)
	authGroup.PUT("/notifications/:id/read", s.Notification.MarkRead)
	authGroup.PUT("/notifications/read-all", s.Notification.MarkAllRead)
}
