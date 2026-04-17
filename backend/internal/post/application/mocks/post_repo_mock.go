package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/blog/blog-community/internal/post/domain"
)

// MockPostRepository is a mock implementation of domain.PostRepository.
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, p *domain.Post) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPostRepository) GetByID(ctx context.Context, id uint64) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*domain.Post), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPostRepository) Update(ctx context.Context, p *domain.Post) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Post, int64, error) {
	args := m.Called(ctx, filter)
	if v := args.Get(0); v != nil {
		return v.([]*domain.Post), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}
