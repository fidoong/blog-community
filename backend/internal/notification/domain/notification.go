package domain

import (
	"context"
	"errors"
	"time"
)

var ErrNotificationNotFound = errors.New("notification not found")

// NotificationType defines notification types.
type NotificationType string

const (
	TypeComment     NotificationType = "comment"
	TypeReply       NotificationType = "reply"
	TypeLikePost    NotificationType = "like_post"
	TypeLikeComment NotificationType = "like_comment"
	TypeFollow      NotificationType = "follow"
	TypeSystem      NotificationType = "system"
)

// Notification represents a user notification.
type Notification struct {
	ID         uint64
	UserID     uint64           // receiver
	Type       NotificationType // comment, reply, like_post, like_comment, follow, system
	Title      string
	Content    string
	ActorID    *uint64
	TargetID   *uint64
	TargetType *string // post, comment, user
	IsRead     bool
	CreatedAt  time.Time
	ReadAt     *time.Time
}

// Notifier abstracts notification sending.
type Notifier interface {
	Send(ctx context.Context, n *Notification) error
}

// Repository defines data access for notifications.
type Repository interface {
	Create(ctx context.Context, n *Notification) error
	GetByID(ctx context.Context, id uint64) (*Notification, error)
	ListByUser(ctx context.Context, userID uint64, onlyUnread bool, page, pageSize int) ([]*Notification, int64, error)
	MarkRead(ctx context.Context, id, userID uint64) error
	MarkAllRead(ctx context.Context, userID uint64) error
	CountUnread(ctx context.Context, userID uint64) (int64, error)
}
