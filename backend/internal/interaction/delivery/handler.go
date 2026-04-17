package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/internal/interaction/application"
)

// InteractionHandler handles HTTP requests for likes and collects.
type InteractionHandler struct {
	useCase application.UseCase
}

func NewInteractionHandler(useCase application.UseCase) *InteractionHandler {
	return &InteractionHandler{useCase: useCase}
}

type toggleResponse struct {
	IsLiked    bool  `json:"isLiked"`
	LikeCount  int64 `json:"likeCount"`
}

type collectResponse struct {
	IsCollected   bool  `json:"isCollected"`
	CollectCount  int64 `json:"collectCount"`
}

func (h *InteractionHandler) ToggleLike(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	targetType := c.Param("targetType")
	targetIDStr := c.Param("targetId")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	isLiked, count, err := h.useCase.ToggleLike(c.Request.Context(), targetType, targetID, claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toggleResponse{IsLiked: isLiked, LikeCount: count})
}

func (h *InteractionHandler) GetLikeStatus(c *gin.Context) {
	var userID uint64
	if claims, ok := middleware.GetAuthUser(c); ok {
		userID = claims.UserID
	}

	targetType := c.Param("targetType")
	targetIDStr := c.Param("targetId")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	isLiked, count, err := h.useCase.GetLikeStatus(c.Request.Context(), targetType, targetID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toggleResponse{IsLiked: isLiked, LikeCount: count})
}

func (h *InteractionHandler) ToggleCollect(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	targetType := c.Param("targetType")
	targetIDStr := c.Param("targetId")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	isCollected, count, err := h.useCase.ToggleCollect(c.Request.Context(), targetType, targetID, claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, collectResponse{IsCollected: isCollected, CollectCount: count})
}

func (h *InteractionHandler) GetCollectStatus(c *gin.Context) {
	var userID uint64
	if claims, ok := middleware.GetAuthUser(c); ok {
		userID = claims.UserID
	}

	targetType := c.Param("targetType")
	targetIDStr := c.Param("targetId")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	isCollected, count, err := h.useCase.GetCollectStatus(c.Request.Context(), targetType, targetID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, collectResponse{IsCollected: isCollected, CollectCount: count})
}
