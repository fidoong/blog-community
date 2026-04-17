package domain

import (
	"context"
	"errors"
	"time"
)

var ErrAlreadyFollowing = errors.New("already following")
var ErrNotFollowing = errors.New("not following")

// Follow represents a follow relationship.
type Follow struct {
	ID          uint64
	FollowerID  uint64
	FollowingID uint64
	CreatedAt   time.Time
}

// FollowStats represents follow statistics for a user.
type FollowStats struct {
	FollowersCount  int64 `json:"followersCount"`
	FollowingCount  int64 `json:"followingCount"`
}

// Repository defines the data access interface for Follow.
type Repository interface {
	Create(ctx context.Context, followerID, followingID uint64) error
	Delete(ctx context.Context, followerID, followingID uint64) error
	IsFollowing(ctx context.Context, followerID, followingID uint64) (bool, error)
	ListFollowers(ctx context.Context, userID uint64, page, pageSize int) ([]*Follow, int64, error)
	ListFollowing(ctx context.Context, userID uint64, page, pageSize int) ([]*Follow, int64, error)
	CountFollowers(ctx context.Context, userID uint64) (int64, error)
	CountFollowing(ctx context.Context, userID uint64) (int64, error)
}
