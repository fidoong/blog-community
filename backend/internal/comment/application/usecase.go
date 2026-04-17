package application

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/blog/blog-community/pkg/errors"
	"github.com/blog/blog-community/internal/comment/domain"
	notificationDomain "github.com/blog/blog-community/internal/notification/domain"
)

// UseCase defines comment application operations.
type UseCase interface {
	Create(ctx context.Context, postID, authorID uint64, content string, parentID *uint64) (*domain.Comment, error)
	GetByID(ctx context.Context, id uint64) (*domain.Comment, error)
	ListByPost(ctx context.Context, postID uint64, page, pageSize int) ([]*domain.Comment, []*domain.Comment, int64, error)
	Delete(ctx context.Context, id, authorID uint64, role string) error
}

type commentUseCase struct {
	repo     domain.CommentRepository
	notifier notificationDomain.Notifier
}

// NewCommentUseCase creates a new comment usecase.
func NewCommentUseCase(repo domain.CommentRepository, notifier notificationDomain.Notifier) UseCase {
	return &commentUseCase{repo: repo, notifier: notifier}
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

	// Send notification asynchronously
	go uc.sendCommentNotification(context.Background(), c)

	return c, nil
}

func (uc *commentUseCase) sendCommentNotification(ctx context.Context, c *domain.Comment) {
	if uc.notifier == nil {
		return
	}

	var recipientID uint64
	var notifType notificationDomain.NotificationType
	var title, targetType string
	var targetID *uint64

	if c.ParentID != nil {
		// Reply to a comment
		authorID, err := uc.repo.GetCommentAuthorID(ctx, *c.ParentID)
		if err != nil || authorID == c.AuthorID {
			return
		}
		recipientID = authorID
		notifType = notificationDomain.TypeReply
		title = "有人回复了你的评论"
		targetType = "comment"
		targetID = c.ParentID
	} else {
		// Comment on a post
		authorID, err := uc.repo.GetPostAuthorID(ctx, c.PostID)
		if err != nil || authorID == c.AuthorID {
			return
		}
		recipientID = authorID
		notifType = notificationDomain.TypeComment
		title = "有人评论了你的文章"
		targetType = "post"
		pid := c.PostID
		targetID = &pid
	}

	_ = uc.notifier.Send(ctx, &notificationDomain.Notification{
		UserID:     recipientID,
		Type:       notifType,
		Title:      title,
		Content:    fmt.Sprintf("%s", c.Content),
		ActorID:    &c.AuthorID,
		TargetID:   targetID,
		TargetType: &targetType,
	})
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
