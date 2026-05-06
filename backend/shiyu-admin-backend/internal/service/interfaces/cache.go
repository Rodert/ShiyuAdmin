package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/vo"
)

// CacheService defines Redis cache management operations.
type CacheService interface {
	ListDatabases(ctx context.Context) ([]*vo.CacheDatabaseVO, error)
	ListKeys(ctx context.Context, db int, pattern, keyType string, count int64) (*vo.CacheKeyListVO, error)
	GetValue(ctx context.Context, db int, key string) (*vo.CacheValueVO, error)
	DeleteKey(ctx context.Context, db int, key string) error
}
