package infrastructure

import (
	"context"
	"fmt"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/user/domain"
)

type txKey struct{}

// EntTransactor manages transactions using the generated ent client.
type EntTransactor struct {
	client *ent.Client
}

// NewEntTransactor creates a new transactor implementing domain.Transactor.
func NewEntTransactor(client *ent.Client) domain.Transactor {
	return &EntTransactor{client: client}
}

func (t *EntTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(txCtx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("rollback failed: %v (original: %w)", rerr, err)
		}
		return err
	}
	return tx.Commit()
}

// ExtractTx extracts the ent client from context if inside a transaction.
func ExtractTx(ctx context.Context, client *ent.Client) *ent.Client {
	if tx, ok := ctx.Value(txKey{}).(*ent.Tx); ok {
		return tx.Client()
	}
	return client
}
