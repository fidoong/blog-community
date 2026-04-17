package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/blog/blog-community/internal/comment/domain"
)

// MockCommentRepository is a mock implementation of domain.CommentRepository.
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, c *domain.Comment) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id uint64) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*domain.Comment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCommentRepository) ListByPost(ctx context.Context, postID uint64, page, pageSize int) ([]*domain.Comment, int64, error) {
	args := m.Called(ctx, postID, page, pageSize)
	if v := args.Get(0); v != nil {
		return v.([]*domain.Comment), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) ListReplies(ctx context.Context, parentID uint64) ([]*domain.Comment, error) {
	args := m.Called(ctx, parentID)
	if v := args.Get(0); v != nil {
		return v.([]*domain.Comment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCommentRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
