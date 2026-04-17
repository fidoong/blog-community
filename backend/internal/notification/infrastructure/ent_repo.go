package infrastructure

import (
	"context"
	"time"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/notification"
	"github.com/blog/blog-community/internal/notification/domain"
)

// EntNotificationRepo implements domain.Repository using Ent.
type EntNotificationRepo struct {
	client *ent.Client
}

// NewEntNotificationRepo creates a new ent-based notification repository.
func NewEntNotificationRepo(client *ent.Client) domain.Repository {
	return &EntNotificationRepo{client: client}
}

func (r *EntNotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	builder := r.client.Notification.Create().
		SetUserID(n.UserID).
		SetType(notification.Type(string(n.Type))).
		SetTitle(n.Title).
		SetContent(n.Content).
		SetIsRead(n.IsRead)

	if n.ActorID != nil {
		builder.SetActorID(*n.ActorID)
	}
	if n.TargetID != nil {
		builder.SetTargetID(*n.TargetID)
	}
	if n.TargetType != nil {
		builder.SetTargetType(*n.TargetType)
	}

	_, err := builder.Save(ctx)
	return err
}

func (r *EntNotificationRepo) GetByID(ctx context.Context, id uint64) (*domain.Notification, error) {
	n, err := r.client.Notification.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}
	return toDomain(n), nil
}

func (r *EntNotificationRepo) ListByUser(ctx context.Context, userID uint64, onlyUnread bool, page, pageSize int) ([]*domain.Notification, int64, error) {
	q := r.client.Notification.Query().Where(notification.UserIDEQ(userID))
	if onlyUnread {
		q = q.Where(notification.IsReadEQ(false))
	}
	q = q.Order(ent.Desc(notification.FieldCreatedAt))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := q.Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toDomainSlice(items), int64(total), nil
}

func (r *EntNotificationRepo) MarkRead(ctx context.Context, id, userID uint64) error {
	return r.client.Notification.UpdateOneID(id).
		Where(notification.UserIDEQ(userID)).
		SetIsRead(true).
		SetReadAt(time.Now()).
		Exec(ctx)
}

func (r *EntNotificationRepo) MarkAllRead(ctx context.Context, userID uint64) error {
	_, err := r.client.Notification.Update().
		Where(notification.UserIDEQ(userID), notification.IsReadEQ(false)).
		SetIsRead(true).
		SetReadAt(time.Now()).
		Save(ctx)
	return err
}

func (r *EntNotificationRepo) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	c, err := r.client.Notification.Query().
		Where(notification.UserIDEQ(userID), notification.IsReadEQ(false)).
		Count(ctx)
	return int64(c), err
}

func toDomainSlice(items []*ent.Notification) []*domain.Notification {
	result := make([]*domain.Notification, len(items))
	for i, n := range items {
		result[i] = toDomain(n)
	}
	return result
}

func toDomain(n *ent.Notification) *domain.Notification {
	var actorID, targetID *uint64
	var targetType *string
	if n.ActorID != nil {
		v := *n.ActorID
		actorID = &v
	}
	if n.TargetID != nil {
		v := *n.TargetID
		targetID = &v
	}
	if n.TargetType != nil {
		v := *n.TargetType
		targetType = &v
	}
	var readAt *time.Time
	if n.ReadAt != nil {
		v := *n.ReadAt
		readAt = &v
	}
	return &domain.Notification{
		ID:         n.ID,
		UserID:     n.UserID,
		Type:       domain.NotificationType(n.Type),
		Title:      n.Title,
		Content:    n.Content,
		ActorID:    actorID,
		TargetID:   targetID,
		TargetType: targetType,
		IsRead:     n.IsRead,
		CreatedAt:  n.CreatedAt,
		ReadAt:     readAt,
	}
}
