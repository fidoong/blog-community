package domain

import (
	"context"
	"errors"
	"time"
)

var ErrPostNotFound = errors.New("post not found")

// Post represents a blog post.
type Post struct {
	ID           uint64
	Title        string
	Content      string
	Summary      string
	ContentType  string
	CoverImage   string
	AuthorID     uint64
	Status       string
	ViewCount    int
	LikeCount    int
	CommentCount int
	CollectCount int
	Tags         []string
	PublishedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const (
	ContentTypeMarkdown = "markdown"
	ContentTypeRichText = "rich_text"

	StatusDraft     = "draft"
	StatusPending   = "pending"
	StatusPublished = "published"
	StatusRejected  = "rejected"
)

// PostRepository defines the data access interface for Post.
type PostRepository interface {
	Create(ctx context.Context, p *Post) error
	GetByID(ctx context.Context, id uint64) (*Post, error)
	Update(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filter ListFilter) ([]*Post, int64, error)
}

// ListFilter defines filtering options for post listing.
type ListFilter struct {
	Status   string
	AuthorID uint64
	Sort     string // "new", "hot"
	Page     int
	PageSize int
}
