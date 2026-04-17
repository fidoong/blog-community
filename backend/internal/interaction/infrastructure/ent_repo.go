package infrastructure

import (
	"context"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/collectrecord"
	"github.com/blog/blog-community/internal/ent/likerecord"
	"github.com/blog/blog-community/internal/interaction/domain"
)

type entInteractionRepo struct {
	client *ent.Client
}

// NewEntInteractionRepo creates a new ent-based interaction repository.
func NewEntInteractionRepo(client *ent.Client) domain.Repository {
	return &entInteractionRepo{client: client}
}

func (r *entInteractionRepo) CreateLike(ctx context.Context, rec *domain.LikeRecord) error {
	created, err := r.client.LikeRecord.Create().
		SetTargetType(likerecord.TargetType(rec.TargetType)).
		SetTargetID(rec.TargetID).
		SetUserID(rec.UserID).
		Save(ctx)
	if err != nil {
		return err
	}
	rec.ID = created.ID
	return nil
}

func (r *entInteractionRepo) DeleteLike(ctx context.Context, targetType string, targetID, userID uint64) error {
	_, err := r.client.LikeRecord.Delete().
		Where(
			likerecord.TargetTypeEQ(likerecord.TargetType(targetType)),
			likerecord.TargetIDEQ(targetID),
			likerecord.UserIDEQ(userID),
		).
		Exec(ctx)
	return err
}

func (r *entInteractionRepo) HasLiked(ctx context.Context, targetType string, targetID, userID uint64) (bool, error) {
	return r.client.LikeRecord.Query().
		Where(
			likerecord.TargetTypeEQ(likerecord.TargetType(targetType)),
			likerecord.TargetIDEQ(targetID),
			likerecord.UserIDEQ(userID),
		).
		Exist(ctx)
}

func (r *entInteractionRepo) CountLikes(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	cnt, err := r.client.LikeRecord.Query().
		Where(
			likerecord.TargetTypeEQ(likerecord.TargetType(targetType)),
			likerecord.TargetIDEQ(targetID),
		).
		Count(ctx)
	return int64(cnt), err
}

func (r *entInteractionRepo) CreateCollect(ctx context.Context, rec *domain.CollectRecord) error {
	created, err := r.client.CollectRecord.Create().
		SetTargetType(collectrecord.TargetType(rec.TargetType)).
		SetTargetID(rec.TargetID).
		SetUserID(rec.UserID).
		Save(ctx)
	if err != nil {
		return err
	}
	rec.ID = created.ID
	return nil
}

func (r *entInteractionRepo) DeleteCollect(ctx context.Context, targetType string, targetID, userID uint64) error {
	_, err := r.client.CollectRecord.Delete().
		Where(
			collectrecord.TargetTypeEQ(collectrecord.TargetType(targetType)),
			collectrecord.TargetIDEQ(targetID),
			collectrecord.UserIDEQ(userID),
		).
		Exec(ctx)
	return err
}

func (r *entInteractionRepo) HasCollected(ctx context.Context, targetType string, targetID, userID uint64) (bool, error) {
	return r.client.CollectRecord.Query().
		Where(
			collectrecord.TargetTypeEQ(collectrecord.TargetType(targetType)),
			collectrecord.TargetIDEQ(targetID),
			collectrecord.UserIDEQ(userID),
		).
		Exist(ctx)
}

func (r *entInteractionRepo) CountCollects(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	cnt, err := r.client.CollectRecord.Query().
		Where(
			collectrecord.TargetTypeEQ(collectrecord.TargetType(targetType)),
			collectrecord.TargetIDEQ(targetID),
		).
		Count(ctx)
	return int64(cnt), err
}
