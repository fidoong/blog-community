package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenStore manages refresh tokens.
type TokenStore interface {
	SaveRefreshToken(ctx context.Context, userID uint64, token string, ttl time.Duration) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (uint64, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

// RedisTokenStore implements TokenStore with Redis.
type RedisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(client *redis.Client) TokenStore {
	return &RedisTokenStore{client: client}
}

func (s *RedisTokenStore) SaveRefreshToken(ctx context.Context, userID uint64, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s", token)
	return s.client.Set(ctx, key, userID, ttl).Err()
}

func (s *RedisTokenStore) GetUserIDByRefreshToken(ctx context.Context, token string) (uint64, error) {
	key := fmt.Sprintf("refresh:%s", token)
	val, err := s.client.Get(ctx, key).Uint64()
	if err == redis.Nil {
		return 0, fmt.Errorf("token not found")
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *RedisTokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh:%s", token)
	return s.client.Del(ctx, key).Err()
}

// GenerateRefreshToken creates a random refresh token.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
