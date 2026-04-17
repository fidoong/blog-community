package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/internal/follow/application"
	"github.com/blog/blog-community/internal/follow/domain"
)

// FollowHandler handles HTTP requests for follows.
type FollowHandler struct {
	useCase application.UseCase
}

// NewFollowHandler creates a new follow handler.
func NewFollowHandler(useCase application.UseCase) *FollowHandler {
	return &FollowHandler{useCase: useCase}
}

type followStatsResponse struct {
	FollowersCount int64 `json:"followersCount"`
	FollowingCount int64 `json:"followingCount"`
	IsFollowing    bool  `json:"isFollowing,omitempty"`
}

type followItemResponse struct {
	ID          uint64 `json:"id"`
	FollowerID  uint64 `json:"followerId"`
	FollowingID uint64 `json:"followingId"`
	CreatedAt   int64  `json:"createdAt"`
}

type listFollowsResponse struct {
	List       []followItemResponse `json:"list"`
	Pagination response.Pagination  `json:"pagination"`
}

func (h *FollowHandler) Follow(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	followingID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	if err := h.useCase.Follow(c.Request.Context(), claims.UserID, followingID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, gin.H{"message": "success"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	idStr := c.Param("id")
	followingID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	if err := h.useCase.Unfollow(c.Request.Context(), claims.UserID, followingID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, gin.H{"message": "success"})
}

func (h *FollowHandler) GetFollowStats(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	stats, err := h.useCase.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := followStatsResponse{
		FollowersCount: stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
	}

	// If viewer is logged in, check if they follow this user
	if claims, ok := middleware.GetAuthUser(c); ok && claims.UserID != userID {
		isFollowing, _ := h.useCase.IsFollowing(c.Request.Context(), claims.UserID, userID)
		resp.IsFollowing = isFollowing
	}

	response.Success(c.Writer, resp)
}

func (h *FollowHandler) ListFollowers(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.useCase.ListFollowers(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toListResponse(items, total, page, pageSize))
}

func (h *FollowHandler) ListFollowing(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.useCase.ListFollowing(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toListResponse(items, total, page, pageSize))
}

func toListResponse(items []*domain.Follow, total int64, page, pageSize int) listFollowsResponse {
	list := make([]followItemResponse, len(items))
	for i, item := range items {
		list[i] = followItemResponse{
			ID:          item.ID,
			FollowerID:  item.FollowerID,
			FollowingID: item.FollowingID,
			CreatedAt:   item.CreatedAt.Unix(),
		}
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return listFollowsResponse{
		List: list,
		Pagination: response.Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}
