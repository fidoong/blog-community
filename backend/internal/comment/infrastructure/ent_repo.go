package infrastructure

import (
	"context"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/comment"
	"github.com/blog/blog-community/internal/comment/domain"
)

type entCommentRepo struct {
	client *ent.Client
}

// NewEntCommentRepo creates a new ent-based comment repository.
func NewEntCommentRepo(client *ent.Client) domain.CommentRepository {
	return &entCommentRepo{client: client}
}

func (r *entCommentRepo) Create(ctx context.Context, c *domain.Comment) error {
	b := r.client.Comment.Create().
		SetContent(c.Content).
		SetPostID(c.PostID).
		SetAuthorID(c.AuthorID)
	if c.ParentID != nil {
		b.SetParentID(*c.ParentID)
	}
	created, err := b.Save(ctx)
	if err != nil {
		return err
	}
	c.ID = created.ID
	return nil
}

func (r *entCommentRepo) GetByID(ctx context.Context, id uint64) (*domain.Comment, error) {
	ec, err := r.client.Comment.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, domain.ErrCommentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(ec), nil
}

func (r *entCommentRepo) ListByPost(ctx context.Context, postID uint64, page, pageSize int) ([]*domain.Comment, int64, error) {
	q := r.client.Comment.Query().
		Where(comment.PostIDEQ(postID), comment.ParentIDIsNil()).
		Order(ent.Desc(comment.FieldCreatedAt))

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	ecs, err := q.Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	comments := make([]*domain.Comment, len(ecs))
	for i, ec := range ecs {
		comments[i] = toDomain(ec)
	}
	return comments, int64(total), nil
}

func (r *entCommentRepo) ListReplies(ctx context.Context, parentID uint64) ([]*domain.Comment, error) {
	ecs, err := r.client.Comment.Query().
		Where(comment.ParentIDEQ(parentID)).
		Order(ent.Asc(comment.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	comments := make([]*domain.Comment, len(ecs))
	for i, ec := range ecs {
		comments[i] = toDomain(ec)
	}
	return comments, nil
}

func (r *entCommentRepo) Delete(ctx context.Context, id uint64) error {
	return r.client.Comment.DeleteOneID(id).Exec(ctx)
}

func (r *entCommentRepo) GetPostAuthorID(ctx context.Context, postID uint64) (uint64, error) {
	p, err := r.client.Post.Get(ctx, postID)
	if err != nil {
		return 0, err
	}
	return p.AuthorID, nil
}

func (r *entCommentRepo) GetCommentAuthorID(ctx context.Context, commentID uint64) (uint64, error) {
	c, err := r.client.Comment.Get(ctx, commentID)
	if err != nil {
		return 0, err
	}
	return c.AuthorID, nil
}

func toDomain(ec *ent.Comment) *domain.Comment {
	c := &domain.Comment{
		ID:        ec.ID,
		Content:   ec.Content,
		PostID:    ec.PostID,
		AuthorID:  ec.AuthorID,
		LikeCount: ec.LikeCount,
		CreatedAt: ec.CreatedAt,
		UpdatedAt: ec.UpdatedAt,
	}
	if ec.ParentID != nil {
		c.ParentID = ec.ParentID
	}
	return c
}
