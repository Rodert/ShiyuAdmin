package cache

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"shiyu-admin-backend/internal/model/vo"
	cacheinterfaces "shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/redis"
)

const redisDatabaseCount = 16

// Service implements CacheService.
type Service struct {
	redisClient *redis.Client
}

// New creates a new CacheService.
func New(redisClient *redis.Client) cacheinterfaces.CacheService {
	return &Service{redisClient: redisClient}
}

// ListDatabases lists the 16 default Redis logical databases and their key counts.
func (s *Service) ListDatabases(ctx context.Context) ([]*vo.CacheDatabaseVO, error) {
	if s.redisClient == nil {
		return []*vo.CacheDatabaseVO{}, nil
	}

	items := make([]*vo.CacheDatabaseVO, 0, redisDatabaseCount)
	for db := 0; db < redisDatabaseCount; db++ {
		count, err := s.redisClient.DBSizeForDB(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("query redis db %d size failed: %w", db, err)
		}
		items = append(items, &vo.CacheDatabaseVO{
			DB:       db,
			KeyCount: count,
		})
	}
	return items, nil
}

// ListKeys scans keys in a Redis logical database.
func (s *Service) ListKeys(ctx context.Context, db int, pattern, keyType string, count int64) (*vo.CacheKeyListVO, error) {
	if s.redisClient == nil {
		return &vo.CacheKeyListVO{DB: db, Pattern: pattern, Type: keyType, Items: []*vo.CacheKeyVO{}}, nil
	}
	if err := validateDB(db); err != nil {
		return nil, err
	}
	if pattern == "" {
		pattern = "*"
	}
	keyType = strings.TrimSpace(strings.ToLower(keyType))
	if keyType == "" {
		keyType = "all"
	}
	if count <= 0 {
		count = 200
	}
	if count > 1000 {
		count = 1000
	}

	keys, err := s.redisClient.ScanKeysForDB(ctx, db, pattern, count)
	if err != nil {
		return nil, err
	}
	sort.Strings(keys)

	items := make([]*vo.CacheKeyVO, 0, len(keys))
	for _, key := range keys {
		actualType, err := s.redisClient.TypeForDB(ctx, db, key)
		if err != nil {
			return nil, err
		}
		if keyType != "all" && actualType != keyType {
			continue
		}
		ttl, err := s.redisClient.TTLForDB(ctx, db, key)
		if err != nil {
			return nil, err
		}
		items = append(items, &vo.CacheKeyVO{
			Key:  key,
			Type: actualType,
			TTL:  ttlSeconds(ttl),
		})
	}

	return &vo.CacheKeyListVO{
		DB:      db,
		Pattern: pattern,
		Type:    keyType,
		Total:   len(items),
		Items:   items,
	}, nil
}

// GetValue returns a Redis key value with type-aware decoding.
func (s *Service) GetValue(ctx context.Context, db int, key string) (*vo.CacheValueVO, error) {
	if s.redisClient == nil {
		return nil, fmt.Errorf("redis is not enabled")
	}
	if err := validateDB(db); err != nil {
		return nil, err
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}

	keyType, value, err := s.redisClient.ValueForDB(ctx, db, key)
	if err != nil {
		return nil, err
	}
	ttl, err := s.redisClient.TTLForDB(ctx, db, key)
	if err != nil {
		return nil, err
	}

	return &vo.CacheValueVO{
		DB:    db,
		Key:   key,
		Type:  keyType,
		TTL:   ttlSeconds(ttl),
		Value: value,
	}, nil
}

func validateDB(db int) error {
	if db < 0 || db >= redisDatabaseCount {
		return fmt.Errorf("redis db must be between 0 and 15")
	}
	return nil
}

func ttlSeconds(ttl time.Duration) int64 {
	if ttl == -1*time.Second {
		return -1
	}
	if ttl == -2*time.Second {
		return -2
	}
	return int64(ttl.Seconds())
}
