package delivery

import (
	"context"
	"fmt"
	"strconv"
	stderrors "errors"
	"time"

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

// ReindexSearch proxies to usecase for startup reindexing.
func (h *PostHandler) ReindexSearch(ctx context.Context) error {
	return h.useCase.ReindexSearch(ctx)
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
	authorIDStr := c.Query("authorId")
	keyword := c.Query("q")

	filter := domain.ListFilter{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Status:   status,
		Keyword:  keyword,
	}

	if authorIDStr != "" {
		authorID, err := strconv.ParseUint(authorIDStr, 10, 64)
		if err == nil {
			filter.AuthorID = authorID
		}
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

func (h *PostHandler) GetRelated(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	posts, err := h.useCase.GetRelated(c.Request.Context(), id, limit)
	if err != nil {
		if stderrors.Is(err, domain.ErrPostNotFound) {
			c.Error(errors.ErrNotFound)
			return
		}
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}

	list := make([]postResponse, len(posts))
	for i, p := range posts {
		list[i] = toPostResponse(p, false)
	}
	response.Success(c.Writer, gin.H{"list": list})
}

func (h *PostHandler) HotKeywords(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	keywords, err := h.useCase.HotKeywords(c.Request.Context(), limit)
	if err != nil {
		c.Error(errors.Wrap(err, errors.ErrInternal))
		return
	}
	response.Success(c.Writer, gin.H{"list": keywords})
}

func (h *PostHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if keyword == "" {
		c.Error(errors.ErrInvalidInput)
		return
	}

	result, err := h.useCase.Search(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	type searchHitResponse struct {
		ID         uint64            `json:"id"`
		Title      string            `json:"title"`
		Summary    string            `json:"summary"`
		Content    string            `json:"content,omitempty"`
		AuthorID   uint64            `json:"authorId"`
		AuthorName string            `json:"authorName"`
		Tags       []string          `json:"tags"`
		ViewCount  int               `json:"viewCount"`
		LikeCount  int               `json:"likeCount"`
		CommentCount int             `json:"commentCount"`
		CreatedAt  int64             `json:"createdAt"`
		Highlight  map[string][]string `json:"highlight,omitempty"`
		Score      float64           `json:"score"`
	}

	list := make([]searchHitResponse, len(result.Hits))
	for i, hit := range result.Hits {
		var id uint64
		fmt.Sscanf(hit.ID, "%d", &id)
		
		s := hit.Source
		r := searchHitResponse{
			ID:         id,
			Score:      hit.Score,
			Highlight:  hit.Highlight,
		}
		if v, ok := s["title"].(string); ok {
			r.Title = v
		}
		if v, ok := s["summary"].(string); ok {
			r.Summary = v
		}
		if v, ok := s["content"].(string); ok {
			r.Content = v
		}
		if v, ok := s["author_name"].(string); ok {
			r.AuthorName = v
		}
		if v, ok := s["author_id"].(float64); ok {
			r.AuthorID = uint64(v)
		}
		if v, ok := s["view_count"].(float64); ok {
			r.ViewCount = int(v)
		}
		if v, ok := s["like_count"].(float64); ok {
			r.LikeCount = int(v)
		}
		if v, ok := s["comment_count"].(float64); ok {
			r.CommentCount = int(v)
		}
		if v, ok := s["tags"].([]any); ok {
			tags := make([]string, len(v))
			for j, t := range v {
				if ts, ok := t.(string); ok {
					tags[j] = ts
				}
			}
			r.Tags = tags
		}
		if v, ok := s["created_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				r.CreatedAt = t.Unix()
			}
		}
		list[i] = r
	}

	totalPages := int(result.Total) / pageSize
	if int(result.Total)%pageSize > 0 {
		totalPages++
	}

	response.Success(c.Writer, gin.H{
		"list": list,
		"pagination": response.Pagination{
			Total:      result.Total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
		"took": result.TookMillis,
	})
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
