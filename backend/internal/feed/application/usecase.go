package application

import (
	"context"

	postdomain "github.com/blog/blog-community/internal/post/domain"
	"github.com/blog/blog-community/internal/feed/domain"
)

// feedUseCase implements domain.UseCase.
type feedUseCase struct {
	postLister domain.PostLister
}

// NewFeedUseCase creates a new feed usecase.
func NewFeedUseCase(postLister domain.PostLister) domain.UseCase {
	return &feedUseCase{postLister: postLister}
}

func (uc *feedUseCase) GetFeed(ctx context.Context, filter domain.FeedFilter) ([]*domain.FeedItem, int64, error) {
	// Normalize pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	switch filter.Type {
	case domain.FeedTypeHot:
		return uc.getHotFeed(ctx, filter)
	case domain.FeedTypeFollowing:
		// TODO: implement following feed after follow system is built
		return []*domain.FeedItem{}, 0, nil
	case domain.FeedTypeRecommend:
		// TODO: implement recommend feed after recommendation algorithm is ready
		return []*domain.FeedItem{}, 0, nil
	default:
		return uc.getLatestFeed(ctx, filter)
	}
}

func (uc *feedUseCase) getLatestFeed(ctx context.Context, filter domain.FeedFilter) ([]*domain.FeedItem, int64, error) {
	posts, total, err := uc.postLister.List(ctx, postdomain.ListFilter{
		Status:   postdomain.StatusPublished,
		Sort:     "new",
		Page:     filter.Page,
		PageSize: filter.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	return toFeedItems(posts), total, nil
}

func (uc *feedUseCase) getHotFeed(ctx context.Context, filter domain.FeedFilter) ([]*domain.FeedItem, int64, error) {
	// For now, hot feed uses post repo's hot sort.
	// Future: filter by period using Redis ZSet.
	posts, total, err := uc.postLister.List(ctx, postdomain.ListFilter{
		Status:   postdomain.StatusPublished,
		Sort:     "hot",
		Page:     filter.Page,
		PageSize: filter.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	return toFeedItems(posts), total, nil
}

func toFeedItems(posts []*postdomain.Post) []*domain.FeedItem {
	items := make([]*domain.FeedItem, len(posts))
	for i, p := range posts {
		items[i] = domain.ToFeedItem(p)
	}
	return items
}
