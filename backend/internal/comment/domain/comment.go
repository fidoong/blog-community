package domain

import (
	"context"
	"errors"
	"time"
)

var ErrCommentNotFound = errors.New("comment not found")

// Comment represents a comment on a post.
type Comment struct {
	ID        uint64
	Content   string
	PostID    uint64
	AuthorID  uint64
	ParentID  *uint64
	LikeCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CommentRepository defines the data access interface for Comment.
type CommentRepository interface {
	Create(ctx context.Context, c *Comment) error
	GetByID(ctx context.Context, id uint64) (*Comment, error)
	ListByPost(ctx context.Context, postID uint64, page, pageSize int) ([]*Comment, int64, error)
	ListReplies(ctx context.Context, parentID uint64) ([]*Comment, error)
	Delete(ctx context.Context, id uint64) error
	GetPostAuthorID(ctx context.Context, postID uint64) (uint64, error)
	GetCommentAuthorID(ctx context.Context, commentID uint64) (uint64, error)
}
