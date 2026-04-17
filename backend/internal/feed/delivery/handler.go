package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/internal/feed/domain"
	"github.com/blog/blog-community/pkg/middleware"
)

// FeedHandler handles HTTP requests for feed.
type FeedHandler struct {
	useCase domain.UseCase
}

// NewFeedHandler creates a new feed handler.
func NewFeedHandler(useCase domain.UseCase) *FeedHandler {
	return &FeedHandler{useCase: useCase}
}

type feedItemResponse struct {
	ID           uint64   `json:"id"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	CoverImage   string   `json:"coverImage"`
	AuthorID     uint64   `json:"authorId"`
	AuthorName   string   `json:"authorName"`
	Tags         []string `json:"tags"`
	LikeCount    int      `json:"likeCount"`
	CommentCount int      `json:"commentCount"`
	ViewCount    int      `json:"viewCount"`
	PublishedAt  int64    `json:"publishedAt"`
	CreatedAt    int64    `json:"createdAt"`
}

type listFeedResponse struct {
	List       []feedItemResponse  `json:"list"`
	Pagination response.Pagination `json:"pagination"`
}

func (h *FeedHandler) GetFeed(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	feedType := c.DefaultQuery("type", "latest")
	period := c.Query("period")

	filter := domain.FeedFilter{
		Type:     domain.FeedType(feedType),
		Page:     page,
		PageSize: pageSize,
		Period:   period,
	}

	// Extract user ID from JWT for following feed
	if claims, ok := middleware.GetAuthUser(c); ok {
		filter.UserID = claims.UserID
	}

	items, total, err := h.useCase.GetFeed(c.Request.Context(), filter)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	list := make([]feedItemResponse, len(items))
	for i, item := range items {
		list[i] = toFeedItemResponse(item)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.Success(c.Writer, listFeedResponse{
		List: list,
		Pagination: response.Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

func toFeedItemResponse(item *domain.FeedItem) feedItemResponse {
	return feedItemResponse{
		ID:           item.ID,
		Title:        item.Title,
		Summary:      item.Summary,
		CoverImage:   item.CoverImage,
		AuthorID:     item.AuthorID,
		AuthorName:   item.AuthorName,
		Tags:         item.Tags,
		LikeCount:    item.LikeCount,
		CommentCount: item.CommentCount,
		ViewCount:    item.ViewCount,
		PublishedAt:  item.PublishedAt,
		CreatedAt:    item.CreatedAt,
	}
}
