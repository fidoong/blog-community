package domain

import (
	"context"
	"time"

	postdomain "github.com/blog/blog-community/internal/post/domain"
)

// FeedItem represents an item in the feed.
type FeedItem struct {
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

// PostLister abstracts post listing for feed generation.
type PostLister interface {
	List(ctx context.Context, filter postdomain.ListFilter) ([]*postdomain.Post, int64, error)
}

// FeedType defines supported feed types.
type FeedType string

const (
	FeedTypeLatest    FeedType = "latest"
	FeedTypeHot       FeedType = "hot"
	FeedTypeFollowing FeedType = "following"
	FeedTypeRecommend FeedType = "recommend"
)

// FeedFilter defines feed query parameters.
type FeedFilter struct {
	Type     FeedType
	Page     int
	PageSize int
	Period   string // "24h", "7d", "30d" for hot feed
}

// UseCase defines feed application operations.
type UseCase interface {
	GetFeed(ctx context.Context, filter FeedFilter) ([]*FeedItem, int64, error)
}

// ToFeedItem converts a post domain model to feed item.
func ToFeedItem(p *postdomain.Post) *FeedItem {
	item := &FeedItem{
		ID:           p.ID,
		Title:        p.Title,
		Summary:      p.Summary,
		CoverImage:   p.CoverImage,
		AuthorID:     p.AuthorID,
		AuthorName:   p.AuthorName,
		Tags:         p.Tags,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
		CreatedAt:    p.CreatedAt.Unix(),
	}
	if !p.PublishedAt.IsZero() {
		item.PublishedAt = p.PublishedAt.Unix()
	}
	return item
}

// DefaultPeriodStart returns the start time for a given period.
func DefaultPeriodStart(period string) time.Time {
	now := time.Now()
	switch period {
	case "24h":
		return now.Add(-24 * time.Hour)
	case "7d":
		return now.Add(-7 * 24 * time.Hour)
	case "30d":
		return now.Add(-30 * 24 * time.Hour)
	default:
		return now.Add(-7 * 24 * time.Hour)
	}
}
