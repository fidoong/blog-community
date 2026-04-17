package application

import (
	"context"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/interaction/domain"
)

// UseCase defines interaction application operations.
type UseCase interface {
	ToggleLike(ctx context.Context, targetType string, targetID, userID uint64) (isLiked bool, count int64, err error)
	GetLikeStatus(ctx context.Context, targetType string, targetID, userID uint64) (isLiked bool, count int64, err error)

	ToggleCollect(ctx context.Context, targetType string, targetID, userID uint64) (isCollected bool, count int64, err error)
	GetCollectStatus(ctx context.Context, targetType string, targetID, userID uint64) (isCollected bool, count int64, err error)
}

type interactionUseCase struct {
	repo    domain.Repository
	counter domain.Counter
}

// NewInteractionUseCase creates a new interaction usecase.
func NewInteractionUseCase(repo domain.Repository, counter domain.Counter) UseCase {
	return &interactionUseCase{repo: repo, counter: counter}
}

func (uc *interactionUseCase) ToggleLike(ctx context.Context, targetType string, targetID, userID uint64) (bool, int64, error) {
	liked, err := uc.repo.HasLiked(ctx, targetType, targetID, userID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}

	if liked {
		if err := uc.repo.DeleteLike(ctx, targetType, targetID, userID); err != nil {
			return false, 0, errors.Wrap(err, errors.ErrInternal)
		}
		count, err := uc.counter.DecrLike(ctx, targetType, targetID)
		if err != nil {
			return false, 0, errors.Wrap(err, errors.ErrInternal)
		}
		return false, count, nil
	}

	rec := &domain.LikeRecord{
		TargetType: targetType,
		TargetID:   targetID,
		UserID:     userID,
	}
	if err := uc.repo.CreateLike(ctx, rec); err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	count, err := uc.counter.IncrLike(ctx, targetType, targetID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	return true, count, nil
}

func (uc *interactionUseCase) GetLikeStatus(ctx context.Context, targetType string, targetID, userID uint64) (bool, int64, error) {
	liked, err := uc.repo.HasLiked(ctx, targetType, targetID, userID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	count, err := uc.counter.GetLikeCount(ctx, targetType, targetID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	if count == 0 {
		// fallback to db and warm cache
		count, err = uc.repo.CountLikes(ctx, targetType, targetID)
		if err != nil {
			return false, 0, errors.Wrap(err, errors.ErrInternal)
		}
		_ = uc.counter.SetLikeCount(ctx, targetType, targetID, count)
	}
	return liked, count, nil
}

func (uc *interactionUseCase) ToggleCollect(ctx context.Context, targetType string, targetID, userID uint64) (bool, int64, error) {
	collected, err := uc.repo.HasCollected(ctx, targetType, targetID, userID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}

	if collected {
		if err := uc.repo.DeleteCollect(ctx, targetType, targetID, userID); err != nil {
			return false, 0, errors.Wrap(err, errors.ErrInternal)
		}
		count, err := uc.counter.DecrCollect(ctx, targetType, targetID)
		if err != nil {
			return false, 0, errors.Wrap(err, errors.ErrInternal)
		}
		return false, count, nil
	}

	rec := &domain.CollectRecord{
		TargetType: targetType,
		TargetID:   targetID,
		UserID:     userID,
	}
	if err := uc.repo.CreateCollect(ctx, rec); err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	count, err := uc.counter.IncrCollect(ctx, targetType, targetID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	return true, count, nil
}

func (uc *interactionUseCase) GetCollectStatus(ctx context.Context, targetType string, targetID, userID uint64) (bool, int64, error) {
	collected, err := uc.repo.HasCollected(ctx, targetType, targetID, userID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	count, err := uc.counter.GetCollectCount(ctx, targetType, targetID)
	if err != nil {
		return false, 0, errors.Wrap(err, errors.ErrInternal)
	}
	if count == 0 {
		count, err = uc.repo.CountCollects(ctx, targetType, targetID)
		if err != nil {
			return false, 0, errors.Wrap(err, errors.ErrInternal)
		}
		_ = uc.counter.SetCollectCount(ctx, targetType, targetID, count)
	}
	return collected, count, nil
}
