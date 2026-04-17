package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/internal/notification/application"
	"github.com/blog/blog-community/internal/notification/domain"
)

// NotificationHandler handles HTTP requests for notifications.
type NotificationHandler struct {
	useCase application.UseCase
}

// NewNotificationHandler creates a new notification handler.
func NewNotificationHandler(useCase application.UseCase) *NotificationHandler {
	return &NotificationHandler{useCase: useCase}
}

type notificationItemResponse struct {
	ID         uint64 `json:"id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ActorID    *uint64 `json:"actorId,omitempty"`
	TargetID   *uint64 `json:"targetId,omitempty"`
	TargetType *string `json:"targetType,omitempty"`
	IsRead     bool   `json:"isRead"`
	CreatedAt  int64  `json:"createdAt"`
}

type listNotificationsResponse struct {
	List       []notificationItemResponse `json:"list"`
	Pagination response.Pagination        `json:"pagination"`
}

type unreadCountResponse struct {
	Count int64 `json:"count"`
}

func (h *NotificationHandler) List(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	onlyUnread := c.Query("unread") == "1" || c.Query("unread") == "true"

	items, total, err := h.useCase.ListByUser(c.Request.Context(), claims.UserID, onlyUnread, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toListResponse(items, total, page, pageSize))
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	if err := h.useCase.MarkRead(c.Request.Context(), id, claims.UserID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, gin.H{"message": "success"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	if err := h.useCase.MarkAllRead(c.Request.Context(), claims.UserID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, gin.H{"message": "success"})
}

func (h *NotificationHandler) CountUnread(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	count, err := h.useCase.CountUnread(c.Request.Context(), claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, unreadCountResponse{Count: count})
}

func toListResponse(items []*domain.Notification, total int64, page, pageSize int) listNotificationsResponse {
	list := make([]notificationItemResponse, len(items))
	for i, item := range items {
		list[i] = notificationItemResponse{
			ID:         item.ID,
			Type:       string(item.Type),
			Title:      item.Title,
			Content:    item.Content,
			ActorID:    item.ActorID,
			TargetID:   item.TargetID,
			TargetType: item.TargetType,
			IsRead:     item.IsRead,
			CreatedAt:  item.CreatedAt.Unix(),
		}
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return listNotificationsResponse{
		List: list,
		Pagination: response.Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}
