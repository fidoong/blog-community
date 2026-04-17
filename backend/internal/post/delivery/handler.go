package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/pkg/validator"
	"github.com/blog/blog-community/internal/post/application"
	"github.com/blog/blog-community/internal/post/domain"
)

// PostHandler handles HTTP requests for posts.
type PostHandler struct {
	useCase application.UseCase
}

func NewPostHandler(useCase application.UseCase) *PostHandler {
	return &PostHandler{useCase: useCase}
}

type createPostRequest struct {
	Title       string   `json:"title" validate:"required,max=256"`
	Content     string   `json:"content" validate:"required"`
	ContentType string   `json:"contentType" validate:"omitempty,oneof=markdown rich_text"`
	CoverImage  string   `json:"coverImage" validate:"omitempty,max=512,url"`
	Tags        []string `json:"tags"`
}

type updatePostRequest struct {
	Title       string   `json:"title" validate:"required,max=256"`
	Content     string   `json:"content" validate:"required"`
	ContentType string   `json:"contentType" validate:"omitempty,oneof=markdown rich_text"`
	CoverImage  string   `json:"coverImage" validate:"omitempty,max=512,url"`
	Tags        []string `json:"tags"`
}

type postResponse struct {
	ID           uint64   `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content,omitempty"`
	Summary      string   `json:"summary"`
	ContentType  string   `json:"contentType"`
	CoverImage   string   `json:"coverImage"`
	AuthorID     uint64   `json:"authorId"`
	AuthorName   string   `json:"authorName"`
	Status       string   `json:"status"`
	ViewCount    int      `json:"viewCount"`
	LikeCount    int      `json:"likeCount"`
	CommentCount int      `json:"commentCount"`
	CollectCount int      `json:"collectCount"`
	Tags         []string `json:"tags"`
	PublishedAt  *int64   `json:"publishedAt,omitempty"`
	CreatedAt    int64    `json:"createdAt"`
	UpdatedAt    int64    `json:"updatedAt"`
}

type listPostsResponse struct {
	List       []postResponse       `json:"list"`
	Pagination response.Pagination  `json:"pagination"`
}

func (h *PostHandler) Create(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	p, err := h.useCase.Create(c.Request.Context(), claims.UserID, req.Title, req.Content, req.ContentType, req.CoverImage, req.Tags)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toPostResponse(p, true))
}

func (h *PostHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	p, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toPostResponse(p, true))
}

func (h *PostHandler) Update(c *gin.Context) {
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

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	p, err := h.useCase.Update(c.Request.Context(), id, claims.UserID, req.Title, req.Content, req.ContentType, req.CoverImage, req.Tags)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toPostResponse(p, true))
}

func (h *PostHandler) Delete(c *gin.Context) {
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

	if err := h.useCase.Delete(c.Request.Context(), id, claims.UserID, claims.Role); err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, gin.H{"message": "success"})
}

func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sort := c.DefaultQuery("sort", "new")
	status := c.Query("status")

	filter := domain.ListFilter{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Status:   status,
	}

	// Default feed only shows published posts
	if filter.Status == "" {
		filter.Status = domain.StatusPublished
	}

	posts, total, err := h.useCase.List(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	list := make([]postResponse, len(posts))
	for i, p := range posts {
		list[i] = toPostResponse(p, false)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	response.Success(c.Writer, listPostsResponse{
		List: list,
		Pagination: response.Pagination{
			Total:      total,
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			TotalPages: totalPages,
		},
	})
}

func (h *PostHandler) Publish(c *gin.Context) {
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

	p, err := h.useCase.Publish(c.Request.Context(), id, claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toPostResponse(p, true))
}

func toPostResponse(p *domain.Post, withContent bool) postResponse {
	resp := postResponse{
		ID:           p.ID,
		Title:        p.Title,
		Summary:      p.Summary,
		ContentType:  p.ContentType,
		CoverImage:   p.CoverImage,
		AuthorID:     p.AuthorID,
		AuthorName:   p.AuthorName,
		Status:       p.Status,
		ViewCount:    p.ViewCount,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		CollectCount: p.CollectCount,
		Tags:         p.Tags,
		CreatedAt:    p.CreatedAt.Unix(),
		UpdatedAt:    p.UpdatedAt.Unix(),
	}
	if withContent {
		resp.Content = p.Content
	}
	if !p.PublishedAt.IsZero() {
		t := p.PublishedAt.Unix()
		resp.PublishedAt = &t
	}
	return resp
}
