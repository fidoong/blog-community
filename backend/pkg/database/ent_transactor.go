package database

import (
	"context"
	"fmt"

	"github.com/blog/blog-community/internal/ent"
)

type entTxKey struct{}

// EntTransactor manages transactions using the generated ent client.
type EntTransactor struct {
	client *ent.Client
}

// NewEntTransactor creates a new transactor.
func NewEntTransactor(client *ent.Client) *EntTransactor {
	return &EntTransactor{client: client}
}

// WithinTransaction executes fn inside an ent transaction.
// The transaction client can be retrieved by ExtractTx inside fn.
func (t *EntTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	txCtx := context.WithValue(ctx, entTxKey{}, tx)
	if err := fn(txCtx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("rollback failed: %v (original: %w)", rerr, err)
		}
		return err
	}
	return tx.Commit()
}

// ExtractTx extracts the ent transaction client from context if inside a transaction,
// otherwise returns the original client.
func ExtractTx(ctx context.Context, client *ent.Client) *ent.Client {
	if tx, ok := ctx.Value(entTxKey{}).(*ent.Tx); ok {
		return tx.Client()
	}
	return client
}
