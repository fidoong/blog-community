package application

import (
	"context"
	stderrors "errors"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/comment/domain"
)

// UseCase defines comment application operations.
type UseCase interface {
	Create(ctx context.Context, postID, authorID uint64, content string, parentID *uint64) (*domain.Comment, error)
	GetByID(ctx context.Context, id uint64) (*domain.Comment, error)
	ListByPost(ctx context.Context, postID uint64, page, pageSize int) ([]*domain.Comment, []*domain.Comment, int64, error)
	Delete(ctx context.Context, id, authorID uint64, role string) error
}

type commentUseCase struct {
	repo domain.CommentRepository
}

// NewCommentUseCase creates a new comment usecase.
func NewCommentUseCase(repo domain.CommentRepository) UseCase {
	return &commentUseCase{repo: repo}
}

func (uc *commentUseCase) Create(ctx context.Context, postID, authorID uint64, content string, parentID *uint64) (*domain.Comment, error) {
	c := &domain.Comment{
		PostID:   postID,
		AuthorID: authorID,
		Content:  content,
		ParentID: parentID,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return c, nil
}

func (uc *commentUseCase) GetByID(ctx context.Context, id uint64) (*domain.Comment, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrCommentNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return c, nil
}

func (uc *commentUseCase) ListByPost(ctx context.Context, postID uint64, page, pageSize int) ([]*domain.Comment, []*domain.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	comments, total, err := uc.repo.ListByPost(ctx, postID, page, pageSize)
	if err != nil {
		return nil, nil, 0, errors.Wrap(err, errors.ErrInternal)
	}

	// Fetch all replies for top-level comments
	var replies []*domain.Comment
	for _, c := range comments {
		if c.ParentID == nil {
			rs, err := uc.repo.ListReplies(ctx, c.ID)
			if err != nil {
				return nil, nil, 0, errors.Wrap(err, errors.ErrInternal)
			}
			replies = append(replies, rs...)
		}
	}
	return comments, replies, total, nil
}

func (uc *commentUseCase) Delete(ctx context.Context, id, authorID uint64, role string) error {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrCommentNotFound) {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	if c.AuthorID != authorID && role != "admin" {
		return errors.ErrForbidden
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}
