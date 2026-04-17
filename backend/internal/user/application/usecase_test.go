package application

import (
	"context"
	"errors"
	"testing"

	"github.com/blog/blog-community/internal/user/application/mocks"
	"github.com/blog/blog-community/internal/user/domain"
	apperrors "github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/pkg/hashutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUseCase_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
			Run(func(args mock.Arguments) {
				u := args.Get(1).(*domain.User)
				u.ID = 1
			}).
			Return(nil)

		u, err := uc.Register(ctx, "test@example.com", "testuser", "Password123!")

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, uint64(1), u.ID)
		assert.Equal(t, "test@example.com", u.Email)
		assert.Equal(t, "testuser", u.Username)
		assert.NotEmpty(t, u.PasswordHash)
		assert.Equal(t, domain.OAuthProviderNone, u.OAuthProvider)
		assert.Empty(t, u.OAuthID)
		assert.Equal(t, domain.RoleUser, u.Role)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
			Return(domain.ErrEmailAlreadyExist)

		u, err := uc.Register(ctx, "dup@example.com", "dupuser", "Password123!")

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, domain.ErrEmailAlreadyExist, err)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate username", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
			Return(domain.ErrUsernameAlreadyExist)

		u, err := uc.Register(ctx, "dup-name@example.com", "dupuser", "Password123!")

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, domain.ErrUsernameAlreadyExist, err)
		repo.AssertExpectations(t)
	})
}

func TestUserUseCase_Login(t *testing.T) {
	ctx := context.Background()
	hash, _ := generateTestHash("Password123!")

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByEmail", ctx, "test@example.com").
			Return(&domain.User{
				ID:           1,
				Email:        "test@example.com",
				Username:     "testuser",
				PasswordHash: hash,
			}, nil)

		u, err := uc.Login(ctx, "test@example.com", "Password123!")

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, uint64(1), u.ID)
		repo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByEmail", ctx, "notfound@example.com").
			Return(nil, domain.ErrUserNotFound)

		u, err := uc.Login(ctx, "notfound@example.com", "Password123!")

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, apperrors.ErrUnauthorized, err)
		repo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByEmail", ctx, "test@example.com").
			Return(&domain.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hash,
			}, nil)

		u, err := uc.Login(ctx, "test@example.com", "WrongPassword!")

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, domain.ErrInvalidPassword, err)
		repo.AssertExpectations(t)
	})
}

func TestUserUseCase_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByID", ctx, uint64(1)).
			Return(&domain.User{ID: 1, Email: "test@example.com"}, nil)

		u, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), u.ID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByID", ctx, uint64(999)).
			Return(nil, domain.ErrUserNotFound)

		u, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, u)
		assert.Equal(t, apperrors.ErrNotFound, err)
		repo.AssertExpectations(t)
	})
}

func TestUserUseCase_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("Update", ctx, mock.AnythingOfType("*domain.User")).
			Return(nil)

		err := uc.Update(ctx, &domain.User{ID: 1, Username: "newname"})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUserUseCase_OAuthLoginOrRegister(t *testing.T) {
	ctx := context.Background()

	t.Run("existing oauth binding", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByOAuth", ctx, "github", "12345").
			Return(&domain.User{ID: 1, Email: "gh@example.com"}, nil)

		u, err := uc.OAuthLoginOrRegister(ctx, "github", "12345", "gh@example.com", "ghuser", "")

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), u.ID)
		repo.AssertExpectations(t)
	})

	t.Run("existing email - bind oauth", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByOAuth", ctx, "github", "12345").
			Return(nil, domain.ErrUserNotFound)
		repo.On("GetByEmail", ctx, "existing@example.com").
			Return(&domain.User{ID: 2, Email: "existing@example.com", AvatarURL: ""}, nil)
		repo.On("Update", ctx, mock.AnythingOfType("*domain.User")).
			Return(nil)

		u, err := uc.OAuthLoginOrRegister(ctx, "github", "12345", "existing@example.com", "ghuser", "http://avatar.png")

		assert.NoError(t, err)
		assert.Equal(t, uint64(2), u.ID)
		assert.Equal(t, "github", u.OAuthProvider)
		assert.Equal(t, "12345", u.OAuthID)
		assert.Equal(t, "http://avatar.png", u.AvatarURL)
		repo.AssertExpectations(t)
	})

	t.Run("new user - auto create", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByOAuth", ctx, "github", "12345").
			Return(nil, domain.ErrUserNotFound)
		repo.On("GetByEmail", ctx, "new@example.com").
			Return(nil, domain.ErrUserNotFound)
		repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).
			Run(func(args mock.Arguments) {
				u := args.Get(1).(*domain.User)
				u.ID = 3
			}).
			Return(nil)

		u, err := uc.OAuthLoginOrRegister(ctx, "github", "12345", "new@example.com", "newuser", "http://avatar.png")

		assert.NoError(t, err)
		assert.Equal(t, uint64(3), u.ID)
		assert.Equal(t, "github", u.OAuthProvider)
		repo.AssertExpectations(t)
	})

	t.Run("getbyoauth error", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		tx := new(mocks.MockTransactor)
		uc := NewUserUseCase(repo, tx)

		repo.On("GetByOAuth", ctx, "github", "12345").
			Return(nil, errors.New("db error"))

		u, err := uc.OAuthLoginOrRegister(ctx, "github", "12345", "test@example.com", "user", "")

		assert.Error(t, err)
		assert.Nil(t, u)
		repo.AssertExpectations(t)
	})
}

// generateTestHash creates a bcrypt hash for test passwords.
func generateTestHash(password string) (string, error) {
	return hashutil.BcryptHash(password)
}
