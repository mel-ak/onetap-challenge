package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"

	"github.com/go-redis/redis/v8"
)

// RedisClient implements CacheService
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

// RateLimit implements rate limiting
func (r *RedisClient) RateLimit(ctx context.Context, key string, limit int, window int64) error {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		r.client.Expire(ctx, key, time.Duration(window)*time.Second)
	}
	if count > int64(limit) {
		return fmt.Errorf("rate limit exceeded")
	}
	return nil
}

// GetBills retrieves cached bills
func (r *RedisClient) GetBills(ctx context.Context, key string) ([]domain.Bill, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var bills []domain.Bill
	return bills, json.Unmarshal([]byte(data), &bills)
}

// CacheBills caches bills
func (r *RedisClient) CacheBills(ctx context.Context, key string, bills []domain.Bill, ttl int64) error {
	data, err := json.Marshal(bills)
	if err != nil {
		return err
	}
	return r.client.SetEX(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}
