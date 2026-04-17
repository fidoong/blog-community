package delivery

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/middleware"
	"github.com/blog/blog-community/pkg/response"
	"github.com/blog/blog-community/pkg/validator"
	"github.com/blog/blog-community/internal/comment/application"
	"github.com/blog/blog-community/internal/comment/domain"
)

// CommentHandler handles HTTP requests for comments.
type CommentHandler struct {
	useCase application.UseCase
}

func NewCommentHandler(useCase application.UseCase) *CommentHandler {
	return &CommentHandler{useCase: useCase}
}

type createCommentRequest struct {
	Content  string  `json:"content" validate:"required,max=2000"`
	ParentID *uint64 `json:"parentId"`
}

type commentResponse struct {
	ID        uint64           `json:"id"`
	Content   string           `json:"content"`
	AuthorID  uint64           `json:"authorId"`
	LikeCount int              `json:"likeCount"`
	Replies   []commentResponse `json:"replies,omitempty"`
	CreatedAt int64            `json:"createdAt"`
}

type listCommentsResponse struct {
	List       []commentResponse `json:"list"`
	Pagination response.Pagination `json:"pagination"`
}

func (h *CommentHandler) Create(c *gin.Context) {
	claims, ok := middleware.GetAuthUser(c)
	if !ok {
		c.Error(errors.ErrUnauthorized)
		return
	}

	postIDStr := c.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}
	if err := validator.Validate(&req); err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	cm, err := h.useCase.Create(c.Request.Context(), postID, claims.UserID, req.Content, req.ParentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c.Writer, toCommentResponse(cm, nil))
}

func (h *CommentHandler) List(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.Error(errors.ErrInvalidInput)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	comments, replies, total, err := h.useCase.ListByPost(c.Request.Context(), postID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	// Build reply map
	replyMap := make(map[uint64][]commentResponse)
	for _, r := range replies {
		replyMap[*r.ParentID] = append(replyMap[*r.ParentID], toCommentResponse(r, nil))
	}

	list := make([]commentResponse, len(comments))
	for i, cm := range comments {
		list[i] = toCommentResponse(cm, replyMap[cm.ID])
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.Success(c.Writer, listCommentsResponse{
		List: list,
		Pagination: response.Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

func (h *CommentHandler) Delete(c *gin.Context) {
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

func toCommentResponse(cm *domain.Comment, replies []commentResponse) commentResponse {
	resp := commentResponse{
		ID:        cm.ID,
		Content:   cm.Content,
		AuthorID:  cm.AuthorID,
		LikeCount: cm.LikeCount,
		CreatedAt: cm.CreatedAt.Unix(),
	}
	if len(replies) > 0 {
		resp.Replies = replies
	}
	return resp
}
