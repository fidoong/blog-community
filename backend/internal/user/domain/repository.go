package domain

import "context"

// UserRepository defines the data access interface for User.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByOAuth(ctx context.Context, provider, oauthID string) (*User, error)
	Update(ctx context.Context, u *User) error
}

// Transactor manages database transactions.
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
