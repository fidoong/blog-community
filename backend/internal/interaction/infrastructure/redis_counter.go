package infrastructure

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/blog/blog-community/internal/interaction/domain"
)

// RedisCounter implements domain.Counter with Redis.
type RedisCounter struct {
	client *redis.Client
}

func NewRedisCounter(client *redis.Client) domain.Counter {
	return &RedisCounter{client: client}
}

func likeKey(targetType string, targetID uint64) string {
	return fmt.Sprintf("like_count:%s:%d", targetType, targetID)
}

func collectKey(targetType string, targetID uint64) string {
	return fmt.Sprintf("collect_count:%s:%d", targetType, targetID)
}

func (r *RedisCounter) IncrLike(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	return r.client.Incr(ctx, likeKey(targetType, targetID)).Result()
}

func (r *RedisCounter) DecrLike(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	return r.client.Decr(ctx, likeKey(targetType, targetID)).Result()
}

func (r *RedisCounter) GetLikeCount(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	val, err := r.client.Get(ctx, likeKey(targetType, targetID)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *RedisCounter) SetLikeCount(ctx context.Context, targetType string, targetID uint64, count int64) error {
	return r.client.Set(ctx, likeKey(targetType, targetID), count, 0).Err()
}

func (r *RedisCounter) IncrCollect(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	return r.client.Incr(ctx, collectKey(targetType, targetID)).Result()
}

func (r *RedisCounter) DecrCollect(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	return r.client.Decr(ctx, collectKey(targetType, targetID)).Result()
}

func (r *RedisCounter) GetCollectCount(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	val, err := r.client.Get(ctx, collectKey(targetType, targetID)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *RedisCounter) SetCollectCount(ctx context.Context, targetType string, targetID uint64, count int64) error {
	return r.client.Set(ctx, collectKey(targetType, targetID), count, 0).Err()
}
