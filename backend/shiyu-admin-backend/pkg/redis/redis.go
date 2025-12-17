package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"shiyu-admin-backend/internal/config"
)

// Client wraps redis client.
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new redis client.
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Get returns value by key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Set sets key-value with expiration.
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// Delete deletes key.
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// Exists checks if key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.rdb.Exists(ctx, key).Result()
	return count > 0, err
}

// Info returns information and statistics about the server.
// Section can be "server", "memory", "stats", "keyspace" or "all".
func (c *Client) Info(ctx context.Context, section string) (string, error) {
	if section == "" {
		section = "all"
	}
	return c.rdb.Info(ctx, section).Result()
}

// DBSize returns the number of keys in the selected database.
func (c *Client) DBSize(ctx context.Context) (int64, error) {
	return c.rdb.DBSize(ctx).Result()
}

// Keys finds all keys matching the given pattern.
// Note: For monitoring/administration only; KEYS can be slow on large datasets.
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}

// Close closes redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

