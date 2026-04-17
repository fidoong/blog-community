package domain

import (
	"context"
	"errors"
	"time"
)

var ErrAlreadyLiked = errors.New("already liked")
var ErrNotLiked = errors.New("not liked")
var ErrAlreadyCollected = errors.New("already collected")
var ErrNotCollected = errors.New("not collected")

// LikeRecord represents a user like action.
type LikeRecord struct {
	ID         uint64
	TargetType string
	TargetID   uint64
	UserID     uint64
	CreatedAt  time.Time
}

// CollectRecord represents a user collect action.
type CollectRecord struct {
	ID         uint64
	TargetType string
	TargetID   uint64
	UserID     uint64
	CreatedAt  time.Time
}

// Repository defines the data access interface for interactions.
type Repository interface {
	CreateLike(ctx context.Context, r *LikeRecord) error
	DeleteLike(ctx context.Context, targetType string, targetID, userID uint64) error
	HasLiked(ctx context.Context, targetType string, targetID, userID uint64) (bool, error)
	CountLikes(ctx context.Context, targetType string, targetID uint64) (int64, error)

	CreateCollect(ctx context.Context, r *CollectRecord) error
	DeleteCollect(ctx context.Context, targetType string, targetID, userID uint64) error
	HasCollected(ctx context.Context, targetType string, targetID, userID uint64) (bool, error)
	CountCollects(ctx context.Context, targetType string, targetID uint64) (int64, error)

	GetPostAuthorID(ctx context.Context, postID uint64) (uint64, error)
	GetCommentAuthorID(ctx context.Context, commentID uint64) (uint64, error)
}

// Counter defines a caching counter for interactions.
type Counter interface {
	IncrLike(ctx context.Context, targetType string, targetID uint64) (int64, error)
	DecrLike(ctx context.Context, targetType string, targetID uint64) (int64, error)
	GetLikeCount(ctx context.Context, targetType string, targetID uint64) (int64, error)
	SetLikeCount(ctx context.Context, targetType string, targetID uint64, count int64) error

	IncrCollect(ctx context.Context, targetType string, targetID uint64) (int64, error)
	DecrCollect(ctx context.Context, targetType string, targetID uint64) (int64, error)
	GetCollectCount(ctx context.Context, targetType string, targetID uint64) (int64, error)
	SetCollectCount(ctx context.Context, targetType string, targetID uint64, count int64) error
}
