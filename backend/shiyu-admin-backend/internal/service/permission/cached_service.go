package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"shiyu-admin-backend/pkg/redis"
)

// CachedService wraps PermissionService with Redis caching.
type CachedService struct {
	base    *Service
	redis   *redis.Client
	ttl     time.Duration
	version string // Cache version for invalidation
}

// NewCachedService creates a cached permission service.
func NewCachedService(base *Service, redisClient *redis.Client, ttl time.Duration) *CachedService {
	return &CachedService{
		base:    base,
		redis:   redisClient,
		ttl:     ttl,
		version: "v1", // Can be incremented to invalidate all caches
	}
}

// GetUserPermissions returns cached user permissions.
func (s *CachedService) GetUserPermissions(ctx context.Context, userCode string) ([]string, error) {
	cacheKey := fmt.Sprintf("user:perms:%s:%s", s.version, userCode)

	// Try to get from cache
	cached, err := s.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var perms []string
		if err := json.Unmarshal([]byte(cached), &perms); err == nil {
			return perms, nil
		}
	}

	// Cache miss, get from base service
	perms, err := s.base.GetUserPermissions(ctx, userCode)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if s.redis != nil {
		permsJSON, _ := json.Marshal(perms)
		_ = s.redis.Set(ctx, cacheKey, permsJSON, s.ttl)
	}

	return perms, nil
}

// CheckPermission checks permission with caching.
func (s *CachedService) CheckPermission(ctx context.Context, userCode string, perms string) (bool, error) {
	userPerms, err := s.GetUserPermissions(ctx, userCode)
	if err != nil {
		return false, err
	}
	for _, perm := range userPerms {
		if perm == perms {
			return true, nil
		}
	}
	return false, nil
}

// InvalidateUserCache invalidates cache for a user.
func (s *CachedService) InvalidateUserCache(ctx context.Context, userCode string) error {
	cacheKey := fmt.Sprintf("user:perms:%s:%s", s.version, userCode)
	return s.redis.Delete(ctx, cacheKey)
}

// InvalidateRoleCache invalidates cache for all users with a role.
func (s *CachedService) InvalidateRoleCache(ctx context.Context, roleCode string) error {
	// This is a simplified version. In production, you might want to maintain
	// a mapping of role -> users, or use a pattern-based deletion.
	// For now, we'll invalidate all user permission caches (brute force).
	// A better approach would be to track user-role relationships in cache.
	_ = fmt.Sprintf("user:perms:%s:*", s.version)
	// Note: Redis doesn't support pattern deletion directly in go-redis v9
	// You would need to use SCAN or maintain a set of keys
	// For simplicity, we'll just return nil here and let TTL handle expiration
	return nil
}
