package application

import (
	"context"
	stderrors "errors"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/user/domain"
)

func (uc *userUseCase) OAuthLoginOrRegister(ctx context.Context, provider, oauthID, email, username, avatarURL string) (*domain.User, error) {
	// 1. Try to find existing oauth binding
	u, err := uc.repo.GetByOAuth(ctx, provider, oauthID)
	if err == nil {
		return u, nil
	}
	if !stderrors.Is(err, domain.ErrUserNotFound) {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// 2. Try to find by email
	u, err = uc.repo.GetByEmail(ctx, email)
	if err == nil {
		// Bind oauth to existing user
		u.OAuthProvider = provider
		u.OAuthID = oauthID
		if avatarURL != "" && u.AvatarURL == "" {
			u.AvatarURL = avatarURL
		}
		if err := uc.repo.Update(ctx, u); err != nil {
			return nil, errors.Wrap(err, errors.ErrInternal)
		}
		return u, nil
	}
	if !stderrors.Is(err, domain.ErrUserNotFound) {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// 3. Create new user
	u = &domain.User{
		Email:         email,
		Username:      username,
		OAuthProvider: provider,
		OAuthID:       oauthID,
		AvatarURL:     avatarURL,
		Role:          domain.RoleUser,
	}
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return u, nil
}
