package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/vo"
)

// MonitorService defines system monitoring operations (cache, online users).
type MonitorService interface {
	GetCacheStats(ctx context.Context) (*vo.CacheStatsVO, error)
	ListOnlineUsers(ctx context.Context) ([]*vo.OnlineUserVO, error)
	UpdateOnlineUser(ctx context.Context, userCode, username, ip, userAgent string) error
}
