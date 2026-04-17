package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/blog/blog-community/internal/interaction/domain"
)

// MockRepository is a mock implementation of domain.Repository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateLike(ctx context.Context, r *domain.LikeRecord) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRepository) DeleteLike(ctx context.Context, targetType string, targetID, userID uint64) error {
	args := m.Called(ctx, targetType, targetID, userID)
	return args.Error(0)
}

func (m *MockRepository) HasLiked(ctx context.Context, targetType string, targetID, userID uint64) (bool, error) {
	args := m.Called(ctx, targetType, targetID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CountLikes(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CreateCollect(ctx context.Context, r *domain.CollectRecord) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRepository) DeleteCollect(ctx context.Context, targetType string, targetID, userID uint64) error {
	args := m.Called(ctx, targetType, targetID, userID)
	return args.Error(0)
}

func (m *MockRepository) HasCollected(ctx context.Context, targetType string, targetID, userID uint64) (bool, error) {
	args := m.Called(ctx, targetType, targetID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CountCollects(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

// MockCounter is a mock implementation of domain.Counter.
type MockCounter struct {
	mock.Mock
}

func (m *MockCounter) IncrLike(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounter) DecrLike(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounter) GetLikeCount(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounter) SetLikeCount(ctx context.Context, targetType string, targetID uint64, count int64) error {
	args := m.Called(ctx, targetType, targetID, count)
	return args.Error(0)
}

func (m *MockCounter) IncrCollect(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounter) DecrCollect(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounter) GetCollectCount(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounter) SetCollectCount(ctx context.Context, targetType string, targetID uint64, count int64) error {
	args := m.Called(ctx, targetType, targetID, count)
	return args.Error(0)
}
