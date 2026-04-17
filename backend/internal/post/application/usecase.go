package application

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/post/domain"
)

// UseCase defines post application operations.
type UseCase interface {
	Create(ctx context.Context, authorID uint64, title, content, contentType, coverImage string, tags []string) (*domain.Post, error)
	GetByID(ctx context.Context, id uint64) (*domain.Post, error)
	Update(ctx context.Context, id, authorID uint64, title, content, contentType, coverImage string, tags []string) (*domain.Post, error)
	Delete(ctx context.Context, id, authorID uint64, role string) error
	List(ctx context.Context, filter domain.ListFilter) ([]*domain.Post, int64, error)
	Publish(ctx context.Context, id, authorID uint64) (*domain.Post, error)
	GetRelated(ctx context.Context, id uint64, limit int) ([]*domain.Post, error)
}

type postUseCase struct {
	repo domain.PostRepository
}

// NewPostUseCase creates a new post usecase.
func NewPostUseCase(repo domain.PostRepository) UseCase {
	return &postUseCase{repo: repo}
}

func (uc *postUseCase) Create(ctx context.Context, authorID uint64, title, content, contentType, coverImage string, tags []string) (*domain.Post, error) {
	p := &domain.Post{
		Title:       title,
		Content:     content,
		ContentType: contentType,
		CoverImage:  coverImage,
		AuthorID:    authorID,
		Status:      domain.StatusDraft,
		Tags:        tags,
	}
	if contentType == "" {
		p.ContentType = domain.ContentTypeMarkdown
	}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return p, nil
}

func (uc *postUseCase) GetByID(ctx context.Context, id uint64) (*domain.Post, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrPostNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return p, nil
}

func (uc *postUseCase) Update(ctx context.Context, id, authorID uint64, title, content, contentType, coverImage string, tags []string) (*domain.Post, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrPostNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	if p.AuthorID != authorID {
		return nil, errors.ErrForbidden
	}

	p.Title = title
	p.Content = content
	if contentType != "" {
		p.ContentType = contentType
	}
	p.CoverImage = coverImage
	p.Tags = tags
	p.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return p, nil
}

func (uc *postUseCase) Delete(ctx context.Context, id, authorID uint64, role string) error {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrPostNotFound) {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}

	if p.AuthorID != authorID && role != "admin" {
		return errors.ErrForbidden
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (uc *postUseCase) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Post, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return uc.repo.List(ctx, filter)
}

func (uc *postUseCase) Publish(ctx context.Context, id, authorID uint64) (*domain.Post, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, domain.ErrPostNotFound) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	if p.AuthorID != authorID {
		return nil, errors.ErrForbidden
	}

	p.Status = domain.StatusPending
	p.PublishedAt = time.Now()

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return p, nil
}

func (uc *postUseCase) GetRelated(ctx context.Context, id uint64, limit int) ([]*domain.Post, error) {
	if limit < 1 || limit > 20 {
		limit = 5
	}
	return uc.repo.GetRelated(ctx, id, limit)
}
