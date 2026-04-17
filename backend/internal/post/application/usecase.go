package application

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/post/domain"
	"github.com/blog/blog-community/pkg/search"
)

// SearchTrendRecorder abstracts search keyword recording.
type SearchTrendRecorder interface {
	RecordKeyword(ctx context.Context, keyword string) error
	HotKeywords(ctx context.Context, n int64) ([]string, error)
}

// PostIndexer abstracts elasticsearch indexing operations.
type PostIndexer interface {
	EnsureIndex(ctx context.Context) error
	IndexPost(ctx context.Context, postID uint64, authorName string) error
	DeletePost(ctx context.Context, postID uint64) error
	SearchPosts(ctx context.Context, keyword string, page, pageSize int) (*search.SearchResult, error)
	ReindexAll(ctx context.Context) error
}

// UseCase defines post application operations.
type UseCase interface {
	Create(ctx context.Context, authorID uint64, title, content, contentType, coverImage string, tags []string) (*domain.Post, error)
	GetByID(ctx context.Context, id uint64) (*domain.Post, error)
	Update(ctx context.Context, id, authorID uint64, title, content, contentType, coverImage string, tags []string) (*domain.Post, error)
	Delete(ctx context.Context, id, authorID uint64, role string) error
	List(ctx context.Context, filter domain.ListFilter) ([]*domain.Post, int64, error)
	Publish(ctx context.Context, id, authorID uint64) (*domain.Post, error)
	GetRelated(ctx context.Context, id uint64, limit int) ([]*domain.Post, error)
	HotKeywords(ctx context.Context, limit int64) ([]string, error)
	Search(ctx context.Context, keyword string, page, pageSize int) (*search.SearchResult, error)
	EnsureSearchIndex(ctx context.Context) error
	ReindexSearch(ctx context.Context) error
}

type postUseCase struct {
	repo       domain.PostRepository
	trendRepo  SearchTrendRecorder
	indexer    PostIndexer
}

// NewPostUseCase creates a new post usecase.
func NewPostUseCase(repo domain.PostRepository, trendRepo SearchTrendRecorder, indexer PostIndexer) UseCase {
	return &postUseCase{repo: repo, trendRepo: trendRepo, indexer: indexer}
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

	// Sync to ES if published
	if p.Status == domain.StatusPublished && uc.indexer != nil {
		go uc.indexer.IndexPost(context.Background(), p.ID, p.AuthorName)
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

	// Remove from ES
	if uc.indexer != nil {
		go uc.indexer.DeletePost(context.Background(), id)
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
	// Record search keyword asynchronously
	if filter.Keyword != "" {
		go uc.trendRepo.RecordKeyword(context.Background(), filter.Keyword)
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

	p.Status = domain.StatusPublished
	p.PublishedAt = time.Now()

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// Sync to ES after publish
	if uc.indexer != nil {
		go uc.indexer.IndexPost(context.Background(), p.ID, p.AuthorName)
	}

	return p, nil
}

func (uc *postUseCase) GetRelated(ctx context.Context, id uint64, limit int) ([]*domain.Post, error) {
	if limit < 1 || limit > 20 {
		limit = 5
	}
	return uc.repo.GetRelated(ctx, id, limit)
}

func (uc *postUseCase) HotKeywords(ctx context.Context, limit int64) ([]string, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return uc.trendRepo.HotKeywords(ctx, limit)
}

func (uc *postUseCase) Search(ctx context.Context, keyword string, page, pageSize int) (*search.SearchResult, error) {
	if keyword == "" {
		return nil, errors.ErrInvalidInput
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Record keyword asynchronously
	go uc.trendRepo.RecordKeyword(context.Background(), keyword)

	if uc.indexer == nil {
		return nil, errors.Wrap(fmt.Errorf("search indexer not available"), errors.ErrInternal)
	}

	return uc.indexer.SearchPosts(ctx, keyword, page, pageSize)
}

func (uc *postUseCase) EnsureSearchIndex(ctx context.Context) error {
	if uc.indexer == nil {
		return nil
	}
	return uc.indexer.EnsureIndex(ctx)
}

func (uc *postUseCase) ReindexSearch(ctx context.Context) error {
	if uc.indexer == nil {
		return nil
	}
	return uc.indexer.ReindexAll(ctx)
}
