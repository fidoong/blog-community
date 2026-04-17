package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/blog/blog-community/internal/interaction/application/mocks"
)

func TestInteractionUseCase_ToggleLike(t *testing.T) {
	ctx := context.Background()

	t.Run("like when not liked", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasLiked", ctx, "post", uint64(1), uint64(100)).Return(false, nil)
		repo.On("CreateLike", ctx, mock.AnythingOfType("*domain.LikeRecord")).Return(nil)
		counter.On("IncrLike", ctx, "post", uint64(1)).Return(int64(10), nil)

		isLiked, count, err := uc.ToggleLike(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.True(t, isLiked)
		assert.Equal(t, int64(10), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})

	t.Run("unlike when already liked", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasLiked", ctx, "post", uint64(1), uint64(100)).Return(true, nil)
		repo.On("DeleteLike", ctx, "post", uint64(1), uint64(100)).Return(nil)
		counter.On("DecrLike", ctx, "post", uint64(1)).Return(int64(9), nil)

		isLiked, count, err := uc.ToggleLike(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.False(t, isLiked)
		assert.Equal(t, int64(9), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})

	t.Run("hasliked error", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasLiked", ctx, "post", uint64(1), uint64(100)).Return(false, errors.New("db error"))

		isLiked, count, err := uc.ToggleLike(ctx, "post", 1, 100)

		assert.Error(t, err)
		assert.False(t, isLiked)
		assert.Equal(t, int64(0), count)
		repo.AssertExpectations(t)
	})
}

func TestInteractionUseCase_GetLikeStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("from cache", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasLiked", ctx, "post", uint64(1), uint64(100)).Return(true, nil)
		counter.On("GetLikeCount", ctx, "post", uint64(1)).Return(int64(50), nil)

		isLiked, count, err := uc.GetLikeStatus(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.True(t, isLiked)
		assert.Equal(t, int64(50), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})

	t.Run("cache miss fallback to db", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasLiked", ctx, "post", uint64(1), uint64(100)).Return(false, nil)
		counter.On("GetLikeCount", ctx, "post", uint64(1)).Return(int64(0), nil)
		repo.On("CountLikes", ctx, "post", uint64(1)).Return(int64(42), nil)
		counter.On("SetLikeCount", ctx, "post", uint64(1), int64(42)).Return(nil)

		isLiked, count, err := uc.GetLikeStatus(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.False(t, isLiked)
		assert.Equal(t, int64(42), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})
}

func TestInteractionUseCase_ToggleCollect(t *testing.T) {
	ctx := context.Background()

	t.Run("collect when not collected", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasCollected", ctx, "post", uint64(1), uint64(100)).Return(false, nil)
		repo.On("CreateCollect", ctx, mock.AnythingOfType("*domain.CollectRecord")).Return(nil)
		counter.On("IncrCollect", ctx, "post", uint64(1)).Return(int64(5), nil)

		isCollected, count, err := uc.ToggleCollect(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.True(t, isCollected)
		assert.Equal(t, int64(5), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})

	t.Run("uncollect when already collected", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasCollected", ctx, "post", uint64(1), uint64(100)).Return(true, nil)
		repo.On("DeleteCollect", ctx, "post", uint64(1), uint64(100)).Return(nil)
		counter.On("DecrCollect", ctx, "post", uint64(1)).Return(int64(4), nil)

		isCollected, count, err := uc.ToggleCollect(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.False(t, isCollected)
		assert.Equal(t, int64(4), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})
}

func TestInteractionUseCase_GetCollectStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("from cache", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasCollected", ctx, "post", uint64(1), uint64(100)).Return(true, nil)
		counter.On("GetCollectCount", ctx, "post", uint64(1)).Return(int64(20), nil)

		isCollected, count, err := uc.GetCollectStatus(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.True(t, isCollected)
		assert.Equal(t, int64(20), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})

	t.Run("cache miss fallback to db", func(t *testing.T) {
		repo := new(mocks.MockRepository)
		counter := new(mocks.MockCounter)
		uc := NewInteractionUseCase(repo, counter)

		repo.On("HasCollected", ctx, "post", uint64(1), uint64(100)).Return(false, nil)
		counter.On("GetCollectCount", ctx, "post", uint64(1)).Return(int64(0), nil)
		repo.On("CountCollects", ctx, "post", uint64(1)).Return(int64(33), nil)
		counter.On("SetCollectCount", ctx, "post", uint64(1), int64(33)).Return(nil)

		isCollected, count, err := uc.GetCollectStatus(ctx, "post", 1, 100)

		assert.NoError(t, err)
		assert.False(t, isCollected)
		assert.Equal(t, int64(33), count)
		repo.AssertExpectations(t)
		counter.AssertExpectations(t)
	})
}
