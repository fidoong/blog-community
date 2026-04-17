package infrastructure

import (
	"context"
	"strings"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/user"
	"github.com/blog/blog-community/internal/user/domain"
	"github.com/blog/blog-community/pkg/database"
)

type entUserRepo struct {
	client *ent.Client
}

// NewEntUserRepo creates a new ent-based user repository.
func NewEntUserRepo(client *ent.Client) domain.UserRepository {
	return &entUserRepo{client: client}
}

func (r *entUserRepo) Create(ctx context.Context, u *domain.User) error {
	client := database.ExtractTx(ctx, r.client)
	creator := client.User.Create().
		SetEmail(u.Email).
		SetUsername(u.Username).
		SetPasswordHash(u.PasswordHash).
		SetOauthProvider(user.OauthProvider(u.OAuthProvider)).
		SetRole(u.Role)
	if avatarURL := strPtr(u.AvatarURL); avatarURL != nil {
		creator.SetAvatarURL(*avatarURL)
	}
	if oauthID := strPtr(u.OAuthID); oauthID != nil {
		creator.SetOauthID(*oauthID)
	}
	created, err := creator.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			switch {
			case strings.Contains(err.Error(), "users_username_key") || strings.Contains(err.Error(), "username"):
				return domain.ErrUsernameAlreadyExist
			case strings.Contains(err.Error(), "users_email_key") || strings.Contains(err.Error(), "email"):
				return domain.ErrEmailAlreadyExist
			}
		}
		return err
	}
	u.ID = created.ID
	return nil
}

func (r *entUserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	client := database.ExtractTx(ctx, r.client)
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
	client := database.ExtractTx(ctx, r.client)
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
	client := database.ExtractTx(ctx, r.client)
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
	client := database.ExtractTx(ctx, r.client)
	updater := client.User.UpdateOneID(u.ID).
		SetUsername(u.Username).
		SetOauthProvider(user.OauthProvider(u.OAuthProvider))
	if avatarURL := strPtr(u.AvatarURL); avatarURL != nil {
		updater.SetAvatarURL(*avatarURL)
	} else {
		updater.ClearAvatarURL()
	}
	if oauthID := strPtr(u.OAuthID); oauthID != nil {
		updater.SetOauthID(*oauthID)
	} else {
		updater.ClearOauthID()
	}
	return updater.Exec(ctx)
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
