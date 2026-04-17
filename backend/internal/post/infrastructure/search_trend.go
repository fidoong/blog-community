package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const searchTrendKey = "search:trend:1h"
const searchTrendTTL = 1 * time.Hour

// SearchTrendRepo handles search keyword trending with Redis.
type SearchTrendRepo struct {
	client *redis.Client
}

// NewSearchTrendRepo creates a new search trend repository.
func NewSearchTrendRepo(client *redis.Client) *SearchTrendRepo {
	return &SearchTrendRepo{client: client}
}

// RecordKeyword increments the score of a search keyword in Redis ZSet.
func (r *SearchTrendRepo) RecordKeyword(ctx context.Context, keyword string) error {
	if keyword == "" {
		return nil
	}
	pipe := r.client.Pipeline()
	pipe.ZIncrBy(ctx, searchTrendKey, 1, keyword)
	pipe.Expire(ctx, searchTrendKey, searchTrendTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// HotKeywords returns the top N trending keywords.
func (r *SearchTrendRepo) HotKeywords(ctx context.Context, n int64) ([]string, error) {
	result, err := r.client.ZRevRange(ctx, searchTrendKey, 0, n-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hot keywords: %w", err)
	}
	return result, nil
}
