package application

import (
	"context"
	stderrors "errors"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/hashutil"
	"github.com/blog/blog-community/internal/user/domain"
)

// UseCase defines user application operations.
type UseCase interface {
	Register(ctx context.Context, email, username, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
	GetByID(ctx context.Context, id uint64) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	OAuthLoginOrRegister(ctx context.Context, provider, oauthID, email, username, avatarURL string) (*domain.User, error)
}

type userUseCase struct {
	repo       domain.UserRepository
	transactor domain.Transactor
}

// NewUserUseCase creates a new user usecase.
func NewUserUseCase(repo domain.UserRepository, transactor domain.Transactor) UseCase {
	return &userUseCase{
		repo:       repo,
		transactor: transactor,
	}
}

func (uc *userUseCase) Register(ctx context.Context, email, username, password string) (*domain.User, error) {
	hashed, err := hashutil.BcryptHash(password)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	u := &domain.User{
		Email:         email,
		Username:      username,
		PasswordHash:  hashed,
		OAuthProvider: domain.OAuthProviderNone,
		Role:          domain.RoleUser,
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		// TODO: better duplicate key detection
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return u, nil
}

func (uc *userUseCase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	u, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if stderrors.Is(err, domain.ErrUserNotFound) {
			return nil, errors.ErrUnauthorized
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	if err := hashutil.BcryptCheck(password, u.PasswordHash); err != nil {
		return nil, domain.ErrInvalidPassword
	}
	return u, nil
}

func (uc *userUseCase) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrUserNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return u, nil
}

func (uc *userUseCase) Update(ctx context.Context, u *domain.User) error {
	if err := uc.repo.Update(ctx, u); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}
