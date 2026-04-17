package application

import (
	"context"
	stderrors "errors"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/follow/domain"
)

// UseCase defines follow application operations.
type UseCase interface {
	Follow(ctx context.Context, followerID, followingID uint64) error
	Unfollow(ctx context.Context, followerID, followingID uint64) error
	IsFollowing(ctx context.Context, followerID, followingID uint64) (bool, error)
	ListFollowers(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Follow, int64, error)
	ListFollowing(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Follow, int64, error)
	GetStats(ctx context.Context, userID uint64) (*domain.FollowStats, error)
}

type followUseCase struct {
	repo domain.Repository
}

// NewFollowUseCase creates a new follow usecase.
func NewFollowUseCase(repo domain.Repository) UseCase {
	return &followUseCase{repo: repo}
}

func (uc *followUseCase) Follow(ctx context.Context, followerID, followingID uint64) error {
	if followerID == followingID {
		return errors.ErrInvalidInput
	}
	if err := uc.repo.Create(ctx, followerID, followingID); err != nil {
		if stderrors.Is(err, domain.ErrAlreadyFollowing) {
			return errors.ErrInvalidInput
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (uc *followUseCase) Unfollow(ctx context.Context, followerID, followingID uint64) error {
	if err := uc.repo.Delete(ctx, followerID, followingID); err != nil {
		if stderrors.Is(err, domain.ErrNotFollowing) {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (uc *followUseCase) IsFollowing(ctx context.Context, followerID, followingID uint64) (bool, error) {
	return uc.repo.IsFollowing(ctx, followerID, followingID)
}

func (uc *followUseCase) ListFollowers(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Follow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.repo.ListFollowers(ctx, userID, page, pageSize)
}

func (uc *followUseCase) ListFollowing(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Follow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.repo.ListFollowing(ctx, userID, page, pageSize)
}

func (uc *followUseCase) GetStats(ctx context.Context, userID uint64) (*domain.FollowStats, error) {
	followers, err := uc.repo.CountFollowers(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	following, err := uc.repo.CountFollowing(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return &domain.FollowStats{
		FollowersCount: followers,
		FollowingCount: following,
	}, nil
}
