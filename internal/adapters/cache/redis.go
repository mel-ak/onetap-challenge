package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mel-ak/onetap-challenge/internal/domain"

	"github.com/redis/go-redis/v9"
)

// RedisClient implements CacheService
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{client: client}
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
func (r *RedisClient) GetBills(ctx context.Context, key string) ([]*domain.Bill, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var bills []*domain.Bill
	return bills, json.Unmarshal([]byte(data), &bills)
}

// CacheBills caches bills
func (r *RedisClient) CacheBills(ctx context.Context, key string, bills []*domain.Bill, ttl int64) error {
	data, err := json.Marshal(bills)
	if err != nil {
		return err
	}
	return r.client.SetEx(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisClient) Client() *redis.Client {
	return c.client
}

func (r *RedisClient) SetWithTTL(ctx context.Context, key string, data interface{}, ttl int) error {
	return r.client.SetEx(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}
