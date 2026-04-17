package infrastructure

import (
	"context"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/user"
	"github.com/blog/blog-community/internal/user/domain"
)

type entUserRepo struct {
	client *ent.Client
}

// NewEntUserRepo creates a new ent-based user repository.
func NewEntUserRepo(client *ent.Client) domain.UserRepository {
	return &entUserRepo{client: client}
}

func (r *entUserRepo) Create(ctx context.Context, u *domain.User) error {
	client := ExtractTx(ctx, r.client)
	created, err := client.User.Create().
		SetEmail(u.Email).
		SetUsername(u.Username).
		SetPasswordHash(u.PasswordHash).
		SetAvatarURL(u.AvatarURL).
		SetOauthProvider(user.OauthProvider(u.OAuthProvider)).
		SetOauthID(u.OAuthID).
		SetRole(u.Role).
		Save(ctx)
	if err != nil {
		return err
	}
	u.ID = created.ID
	return nil
}

func (r *entUserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	client := ExtractTx(ctx, r.client)
	eu, err := client.User.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(eu), nil
}

func (r *entUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	client := ExtractTx(ctx, r.client)
	eu, err := client.User.Query().Where(user.Email(email)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(eu), nil
}

func (r *entUserRepo) GetByOAuth(ctx context.Context, provider, oauthID string) (*domain.User, error) {
	client := ExtractTx(ctx, r.client)
	eu, err := client.User.Query().
		Where(user.OauthProviderEQ(user.OauthProvider(provider)), user.OauthID(oauthID)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(eu), nil
}

func (r *entUserRepo) Update(ctx context.Context, u *domain.User) error {
	client := ExtractTx(ctx, r.client)
	return client.User.UpdateOneID(u.ID).
		SetUsername(u.Username).
		SetAvatarURL(u.AvatarURL).
		SetOauthID(u.OAuthID).
		SetOauthProvider(user.OauthProvider(u.OAuthProvider)).
		Exec(ctx)
}

func toDomain(eu *ent.User) *domain.User {
	return &domain.User{
		ID:            eu.ID,
		Email:         eu.Email,
		Username:      eu.Username,
		PasswordHash:  eu.PasswordHash,
		AvatarURL:     eu.AvatarURL,
		OAuthProvider: string(eu.OauthProvider),
		OAuthID:       eu.OauthID,
		Role:          eu.Role,
		CreatedAt:     eu.CreatedAt,
		UpdatedAt:     eu.UpdatedAt,
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
