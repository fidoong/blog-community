package infrastructure

import (
	"context"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/post"
	postDomain "github.com/blog/blog-community/internal/post/domain"
)

type entPostRepo struct {
	client *ent.Client
}

// NewEntPostRepo creates a new ent-based post repository.
func NewEntPostRepo(client *ent.Client) postDomain.PostRepository {
	return &entPostRepo{client: client}
}

func (r *entPostRepo) Create(ctx context.Context, p *postDomain.Post) error {
	b := r.client.Post.Create().
		SetTitle(p.Title).
		SetContent(p.Content).
		SetSummary(p.Summary).
		SetContentType(post.ContentType(p.ContentType)).
		SetCoverImage(p.CoverImage).
		SetAuthorID(p.AuthorID).
		SetStatus(post.Status(p.Status)).
		SetTags(p.Tags)
	if !p.PublishedAt.IsZero() {
		b.SetPublishedAt(p.PublishedAt)
	}
	created, err := b.Save(ctx)
	if err != nil {
		return err
	}
	p.ID = created.ID
	return nil
}

func (r *entPostRepo) GetByID(ctx context.Context, id uint64) (*postDomain.Post, error) {
	ep, err := r.client.Post.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, postDomain.ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(ep), nil
}

func (r *entPostRepo) Update(ctx context.Context, p *postDomain.Post) error {
	b := r.client.Post.UpdateOneID(p.ID).
		SetTitle(p.Title).
		SetContent(p.Content).
		SetSummary(p.Summary).
		SetContentType(post.ContentType(p.ContentType)).
		SetCoverImage(p.CoverImage).
		SetStatus(post.Status(p.Status)).
		SetTags(p.Tags)
	if !p.PublishedAt.IsZero() {
		b.SetPublishedAt(p.PublishedAt)
	}
	return b.Exec(ctx)
}

func (r *entPostRepo) Delete(ctx context.Context, id uint64) error {
	return r.client.Post.DeleteOneID(id).Exec(ctx)
}

func (r *entPostRepo) List(ctx context.Context, filter postDomain.ListFilter) ([]*postDomain.Post, int64, error) {
	q := r.client.Post.Query()

	if filter.Status != "" {
		q = q.Where(post.StatusEQ(post.Status(filter.Status)))
	}
	if filter.AuthorID > 0 {
		q = q.Where(post.AuthorIDEQ(filter.AuthorID))
	}

	switch filter.Sort {
	case "hot":
		q = q.Order(ent.Desc(post.FieldLikeCount), ent.Desc(post.FieldCreatedAt))
	default:
		q = q.Order(ent.Desc(post.FieldCreatedAt))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	eps, err := q.Offset((page - 1) * pageSize).Limit(pageSize).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	posts := make([]*postDomain.Post, len(eps))
	for i, ep := range eps {
		posts[i] = toDomain(ep)
	}
	return posts, int64(total), nil
}

func toDomain(ep *ent.Post) *postDomain.Post {
	return &postDomain.Post{
		ID:           ep.ID,
		Title:        ep.Title,
		Content:      ep.Content,
		Summary:      ep.Summary,
		ContentType:  string(ep.ContentType),
		CoverImage:   ep.CoverImage,
		AuthorID:     ep.AuthorID,
		Status:       string(ep.Status),
		ViewCount:    ep.ViewCount,
		LikeCount:    ep.LikeCount,
		CommentCount: ep.CommentCount,
		CollectCount: ep.CollectCount,
		Tags:         ep.Tags,
		PublishedAt:  ep.PublishedAt,
		CreatedAt:    ep.CreatedAt,
		UpdatedAt:    ep.UpdatedAt,
	}
}
