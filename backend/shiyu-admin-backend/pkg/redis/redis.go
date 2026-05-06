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

// DeleteForDB deletes a key in a specific Redis logical database.
func (c *Client) DeleteForDB(ctx context.Context, db int, key string) error {
	client := c.clientForDB(db)
	defer client.Close()
	return client.Del(ctx, key).Err()
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

// DBSizeForDB returns the number of keys in a specific Redis logical database.
func (c *Client) DBSizeForDB(ctx context.Context, db int) (int64, error) {
	client := c.clientForDB(db)
	defer client.Close()
	return client.DBSize(ctx).Result()
}

// Keys finds all keys matching the given pattern.
// Note: For monitoring/administration only; KEYS can be slow on large datasets.
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}

// ScanKeys incrementally scans keys by pattern to avoid blocking Redis.
func (c *Client) ScanKeys(ctx context.Context, pattern string, count int64) ([]string, error) {
	if count <= 0 {
		count = 200
	}

	cursor := uint64(0)
	result := make([]string, 0, 64)
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}
		result = append(result, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

// ScanKeysForDB incrementally scans keys in a specific Redis logical database.
func (c *Client) ScanKeysForDB(ctx context.Context, db int, pattern string, count int64) ([]string, error) {
	if count <= 0 {
		count = 200
	}
	if pattern == "" {
		pattern = "*"
	}

	client := c.clientForDB(db)
	defer client.Close()

	cursor := uint64(0)
	result := make([]string, 0, 64)
	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}
		result = append(result, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

// TypeForDB returns the Redis data type for a key in a specific database.
func (c *Client) TypeForDB(ctx context.Context, db int, key string) (string, error) {
	client := c.clientForDB(db)
	defer client.Close()
	return client.Type(ctx, key).Result()
}

// TTLForDB returns the TTL for a key in a specific database.
func (c *Client) TTLForDB(ctx context.Context, db int, key string) (time.Duration, error) {
	client := c.clientForDB(db)
	defer client.Close()
	return client.TTL(ctx, key).Result()
}

// ValueForDB returns the value for a key in a specific database, normalized by Redis type.
func (c *Client) ValueForDB(ctx context.Context, db int, key string) (string, interface{}, error) {
	client := c.clientForDB(db)
	defer client.Close()

	keyType, err := client.Type(ctx, key).Result()
	if err != nil {
		return "", nil, err
	}

	switch keyType {
	case "string":
		value, err := client.Get(ctx, key).Result()
		return keyType, value, err
	case "list":
		value, err := client.LRange(ctx, key, 0, -1).Result()
		return keyType, value, err
	case "set":
		value, err := client.SMembers(ctx, key).Result()
		return keyType, value, err
	case "zset":
		items, err := client.ZRangeWithScores(ctx, key, 0, -1).Result()
		if err != nil {
			return keyType, nil, err
		}
		value := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			value = append(value, map[string]interface{}{
				"member": item.Member,
				"score":  item.Score,
			})
		}
		return keyType, value, nil
	case "hash":
		value, err := client.HGetAll(ctx, key).Result()
		return keyType, value, err
	case "stream":
		items, err := client.XRangeN(ctx, key, "-", "+", 100).Result()
		if err != nil {
			return keyType, nil, err
		}
		value := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			value = append(value, map[string]interface{}{
				"id":     item.ID,
				"values": item.Values,
			})
		}
		return keyType, value, nil
	default:
		return keyType, nil, nil
	}
}

func (c *Client) clientForDB(db int) *redis.Client {
	opts := *c.rdb.Options()
	opts.DB = db
	return redis.NewClient(&opts)
}

// Close closes redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}
