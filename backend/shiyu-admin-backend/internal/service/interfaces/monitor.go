package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/vo"
)

// MonitorService defines system monitoring operations (cache, online users).
type MonitorService interface {
	GetCacheStats(ctx context.Context) (*vo.CacheStatsVO, error)
	GetDatabaseStats(ctx context.Context) (*vo.DatabaseStatsVO, error)
	ListOnlineUsers(ctx context.Context) ([]*vo.OnlineUserVO, error)
	UpdateOnlineUser(ctx context.Context, sessionID, userCode, username, ip, hostIP, userAgent string) error
	ForceLogout(ctx context.Context, sessionID string) error
	IsSessionRevoked(ctx context.Context, sessionID string) (bool, error)
}
