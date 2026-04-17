package application

import (
	"context"
	stderrors "errors"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/notification/domain"
)

// UseCase defines notification application operations.
type UseCase interface {
	domain.Notifier
	ListByUser(ctx context.Context, userID uint64, onlyUnread bool, page, pageSize int) ([]*domain.Notification, int64, error)
	MarkRead(ctx context.Context, id, userID uint64) error
	MarkAllRead(ctx context.Context, userID uint64) error
	CountUnread(ctx context.Context, userID uint64) (int64, error)
}

type notificationUseCase struct {
	repo domain.Repository
}

// NewNotificationUseCase creates a new notification usecase.
func NewNotificationUseCase(repo domain.Repository) UseCase {
	return &notificationUseCase{repo: repo}
}

func (uc *notificationUseCase) Send(ctx context.Context, n *domain.Notification) error {
	if n.UserID == 0 {
		return errors.ErrInvalidInput
	}
	if n.Type == "" {
		n.Type = domain.TypeSystem
	}
	if err := uc.repo.Create(ctx, n); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (uc *notificationUseCase) ListByUser(ctx context.Context, userID uint64, onlyUnread bool, page, pageSize int) ([]*domain.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.repo.ListByUser(ctx, userID, onlyUnread, page, pageSize)
}

func (uc *notificationUseCase) MarkRead(ctx context.Context, id, userID uint64) error {
	n, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrNotificationNotFound) {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	if n.UserID != userID {
		return errors.ErrForbidden
	}
	if err := uc.repo.MarkRead(ctx, id, userID); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (uc *notificationUseCase) MarkAllRead(ctx context.Context, userID uint64) error {
	if err := uc.repo.MarkAllRead(ctx, userID); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (uc *notificationUseCase) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return uc.repo.CountUnread(ctx, userID)
}
