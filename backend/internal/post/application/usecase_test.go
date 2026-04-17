package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/blog/blog-community/internal/post/application/mocks"
	"github.com/blog/blog-community/internal/post/domain"
	apperrors "github.com/blog/blog-community/pkg/errors"
)

func TestPostUseCase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success with markdown", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.Post")).
			Run(func(args mock.Arguments) {
				p := args.Get(1).(*domain.Post)
				p.ID = 1
			}).
			Return(nil)

		p, err := uc.Create(ctx, 100, "Title", "Content", "", "", []string{"Go"})

		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, uint64(1), p.ID)
		assert.Equal(t, "Title", p.Title)
		assert.Equal(t, "Content", p.Content)
		assert.Equal(t, domain.ContentTypeMarkdown, p.ContentType)
		assert.Equal(t, domain.StatusDraft, p.Status)
		assert.Equal(t, uint64(100), p.AuthorID)
		repo.AssertExpectations(t)
	})

	t.Run("success with rich text", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("Create", ctx, mock.AnythingOfType("*domain.Post")).
			Run(func(args mock.Arguments) {
				p := args.Get(1).(*domain.Post)
				p.ID = 2
			}).
			Return(nil)

		p, err := uc.Create(ctx, 100, "Title", "Content", domain.ContentTypeRichText, "", nil)

		assert.NoError(t, err)
		assert.Equal(t, domain.ContentTypeRichText, p.ContentType)
		repo.AssertExpectations(t)
	})
}

func TestPostUseCase_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		expected := &domain.Post{ID: 1, Title: "Test"}
		repo.On("GetByID", ctx, uint64(1)).Return(expected, nil)

		p, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, p)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("GetByID", ctx, uint64(999)).Return(nil, domain.ErrPostNotFound)

		p, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Equal(t, apperrors.ErrNotFound, err)
		repo.AssertExpectations(t)
	})
}

func TestPostUseCase_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		existing := &domain.Post{ID: 1, Title: "Old", Content: "Old", AuthorID: 100, Status: domain.StatusDraft}
		repo.On("GetByID", ctx, uint64(1)).Return(existing, nil)
		repo.On("Update", ctx, mock.AnythingOfType("*domain.Post")).Return(nil)

		p, err := uc.Update(ctx, 1, 100, "New Title", "New Content", "", "", nil)

		assert.NoError(t, err)
		assert.Equal(t, "New Title", p.Title)
		assert.Equal(t, "New Content", p.Content)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("GetByID", ctx, uint64(999)).Return(nil, domain.ErrPostNotFound)

		p, err := uc.Update(ctx, 999, 100, "Title", "Content", "", "", nil)

		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Equal(t, apperrors.ErrNotFound, err)
		repo.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		existing := &domain.Post{ID: 1, Title: "Old", AuthorID: 100, Status: domain.StatusDraft}
		repo.On("GetByID", ctx, uint64(1)).Return(existing, nil)

		p, err := uc.Update(ctx, 1, 200, "New Title", "Content", "", "", nil)

		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Equal(t, apperrors.ErrForbidden, err)
		repo.AssertExpectations(t)
	})
}

func TestPostUseCase_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Post{ID: 1, AuthorID: 100}, nil)
		repo.On("Delete", ctx, uint64(1)).Return(nil)

		err := uc.Delete(ctx, 1, 100, "user")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("admin can delete others", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Post{ID: 1, AuthorID: 100}, nil)
		repo.On("Delete", ctx, uint64(1)).Return(nil)

		err := uc.Delete(ctx, 1, 999, "admin")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Post{ID: 1, AuthorID: 100}, nil)

		err := uc.Delete(ctx, 1, 200, "user")

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden, err)
		repo.AssertExpectations(t)
	})
}

func TestPostUseCase_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		posts := []*domain.Post{
			{ID: 1, Title: "Post 1", CreatedAt: time.Now()},
			{ID: 2, Title: "Post 2", CreatedAt: time.Now()},
		}
		repo.On("List", ctx, mock.AnythingOfType("domain.ListFilter")).Return(posts, int64(2), nil)

		result, total, err := uc.List(ctx, domain.ListFilter{Page: 1, PageSize: 20, Sort: "new"})

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
		repo.AssertExpectations(t)
	})

	t.Run("default pagination", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("List", ctx, mock.AnythingOfType("domain.ListFilter")).Return([]*domain.Post{}, int64(0), nil)

		_, _, err := uc.List(ctx, domain.ListFilter{})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestPostUseCase_Publish(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		existing := &domain.Post{ID: 1, Title: "Draft", AuthorID: 100, Status: domain.StatusDraft}
		repo.On("GetByID", ctx, uint64(1)).Return(existing, nil)
		repo.On("Update", ctx, mock.AnythingOfType("*domain.Post")).Return(nil)

		p, err := uc.Publish(ctx, 1, 100)

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusPending, p.Status)
		assert.False(t, p.PublishedAt.IsZero())
		repo.AssertExpectations(t)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := new(mocks.MockPostRepository)
		uc := NewPostUseCase(repo)

		repo.On("GetByID", ctx, uint64(1)).Return(&domain.Post{ID: 1, AuthorID: 100}, nil)

		p, err := uc.Publish(ctx, 1, 200)

		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Equal(t, apperrors.ErrForbidden, err)
		repo.AssertExpectations(t)
	})
}
