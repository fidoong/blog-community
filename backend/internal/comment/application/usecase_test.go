package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/blog/blog-community/internal/comment/application/mocks"
	"github.com/blog/blog-community/internal/comment/domain"
	apperrors "github.com/blog/blog-community/pkg/errors"
)

func ptr[T any](v T) *T { return &v }

func TestCommentUseCase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success top-level", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.Comment")).
			Run(func(args mock.Arguments) {
				c := args.Get(1).(*domain.Comment)
				c.ID = 1
			}).
			Return(nil)

		c, err := uc.Create(ctx, 10, 100, "Great post!", nil)

		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, uint64(1), c.ID)
		assert.Equal(t, uint64(10), c.PostID)
		assert.Equal(t, uint64(100), c.AuthorID)
		assert.Equal(t, "Great post!", c.Content)
		assert.Nil(t, c.ParentID)
		repo.AssertExpectations(t)
	})

	t.Run("success reply", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.Comment")).
			Run(func(args mock.Arguments) {
				c := args.Get(1).(*domain.Comment)
				c.ID = 2
			}).
			Return(nil)

		c, err := uc.Create(ctx, 10, 100, "Thanks!", ptr(uint64(1)))

		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, uint64(2), c.ID)
		assert.Equal(t, uint64(1), *c.ParentID)
		repo.AssertExpectations(t)
	})
}

func TestCommentUseCase_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Comment{ID: 1, Content: "Test"}, nil)

		c, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), c.ID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("GetByID", ctx, uint64(999)).Return(nil, domain.ErrCommentNotFound)

		c, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, c)
		assert.Equal(t, apperrors.ErrNotFound, err)
		repo.AssertExpectations(t)
	})
}

func TestCommentUseCase_ListByPost(t *testing.T) {
	ctx := context.Background()

	t.Run("success with replies", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		comments := []*domain.Comment{
			{ID: 1, PostID: 10, AuthorID: 100, Content: "Top 1", ParentID: nil},
			{ID: 2, PostID: 10, AuthorID: 101, Content: "Top 2", ParentID: nil},
		}
		replies := []*domain.Comment{
			{ID: 3, PostID: 10, AuthorID: 102, Content: "Reply 1", ParentID: ptr(uint64(1))},
		}

		repo.On("ListByPost", ctx, uint64(10), 1, 20).Return(comments, int64(2), nil)
		repo.On("ListReplies", ctx, uint64(1)).Return(replies, nil)
		repo.On("ListReplies", ctx, uint64(2)).Return([]*domain.Comment{}, nil)

		result, replyResult, total, err := uc.ListByPost(ctx, 10, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Len(t, replyResult, 1)
		assert.Equal(t, int64(2), total)
		repo.AssertExpectations(t)
	})

	t.Run("default pagination", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("ListByPost", ctx, uint64(10), 1, 20).Return([]*domain.Comment{}, int64(0), nil)

		_, _, _, err := uc.ListByPost(ctx, 10, 0, 0)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestCommentUseCase_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Comment{ID: 1, AuthorID: 100}, nil)
		repo.On("Delete", ctx, uint64(1)).Return(nil)

		err := uc.Delete(ctx, 1, 100, "user")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("admin can delete others", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Comment{ID: 1, AuthorID: 100}, nil)
		repo.On("Delete", ctx, uint64(1)).Return(nil)

		err := uc.Delete(ctx, 1, 999, "admin")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("GetByID", ctx, uint64(999)).Return(nil, domain.ErrCommentNotFound)

		err := uc.Delete(ctx, 999, 100, "user")

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrNotFound, err)
		repo.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := new(mocks.MockCommentRepository)
		uc := NewCommentUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Comment{ID: 1, AuthorID: 100}, nil)

		err := uc.Delete(ctx, 1, 200, "user")

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		repo.AssertExpectations(t)
	})
}
