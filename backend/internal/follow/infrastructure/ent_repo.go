package infrastructure

import (
	"context"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/follow"
	"github.com/blog/blog-community/internal/follow/domain"
)

// EntFollowRepo implements domain.Repository using Ent.
type EntFollowRepo struct {
	client *ent.Client
}

// NewEntFollowRepo creates a new ent-based follow repository.
func NewEntFollowRepo(client *ent.Client) domain.Repository {
	return &EntFollowRepo{client: client}
}

func (r *EntFollowRepo) Create(ctx context.Context, followerID, followingID uint64) error {
	_, err := r.client.Follow.Create().
		SetFollowerID(followerID).
		SetFollowingID(followingID).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return domain.ErrAlreadyFollowing
	}
	return err
}

func (r *EntFollowRepo) Delete(ctx context.Context, followerID, followingID uint64) error {
	n, err := r.client.Follow.Delete().
		Where(follow.FollowerIDEQ(followerID), follow.FollowingIDEQ(followingID)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFollowing
	}
	return nil
}

func (r *EntFollowRepo) IsFollowing(ctx context.Context, followerID, followingID uint64) (bool, error) {
	n, err := r.client.Follow.Query().
		Where(follow.FollowerIDEQ(followerID), follow.FollowingIDEQ(followingID)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *EntFollowRepo) ListFollowers(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Follow, int64, error) {
	q := r.client.Follow.Query().Where(follow.FollowingIDEQ(userID)).Order(ent.Desc(follow.FieldCreatedAt))
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return toDomainSlice(items), int64(total), nil
}

func (r *EntFollowRepo) ListFollowing(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Follow, int64, error) {
	q := r.client.Follow.Query().Where(follow.FollowerIDEQ(userID)).Order(ent.Desc(follow.FieldCreatedAt))
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, err := q.Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return toDomainSlice(items), int64(total), nil
}

func (r *EntFollowRepo) CountFollowers(ctx context.Context, userID uint64) (int64, error) {
	c, err := r.client.Follow.Query().Where(follow.FollowingIDEQ(userID)).Count(ctx)
	return int64(c), err
}

func (r *EntFollowRepo) CountFollowing(ctx context.Context, userID uint64) (int64, error) {
	c, err := r.client.Follow.Query().Where(follow.FollowerIDEQ(userID)).Count(ctx)
	return int64(c), err
}

func (r *EntFollowRepo) ListFollowingIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	items, err := r.client.Follow.Query().
		Where(follow.FollowerIDEQ(userID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, len(items))
	for i, item := range items {
		ids[i] = item.FollowingID
	}
	return ids, nil
}

func toDomainSlice(items []*ent.Follow) []*domain.Follow {
	result := make([]*domain.Follow, len(items))
	for i, f := range items {
		result[i] = toDomain(f)
	}
	return result
}

func toDomain(f *ent.Follow) *domain.Follow {
	return &domain.Follow{
		ID:          f.ID,
		FollowerID:  f.FollowerID,
		FollowingID: f.FollowingID,
		CreatedAt:   f.CreatedAt,
	}
}
