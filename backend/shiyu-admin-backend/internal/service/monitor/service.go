package monitor

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"shiyu-admin-backend/internal/model/vo"
	monitorinterfaces "shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/redis"
)

// Service implements MonitorService.
type Service struct {
	redisClient *redis.Client
	// onlineTTL defines how long a user is considered online since last activity.
	onlineTTL time.Duration
}

// New creates a new MonitorService.
func New(redisClient *redis.Client, onlineTTL time.Duration) monitorinterfaces.MonitorService {
	return &Service{
		redisClient: redisClient,
		onlineTTL:   onlineTTL,
	}
}

// GetCacheStats returns basic Redis statistics for monitoring.
func (s *Service) GetCacheStats(ctx context.Context) (*vo.CacheStatsVO, error) {
	if s.redisClient == nil {
		return nil, nil
	}

	info, err := s.redisClient.Info(ctx, "all")
	if err != nil {
		return nil, err
	}

	parsed := parseRedisInfo(info)

	var stats vo.CacheStatsVO
	stats.RedisVersion = parsed["redis_version"]
	stats.Mode = parsed["redis_mode"]
	stats.UsedMemory, _ = parseInt64(parsed["used_memory"])
	stats.UsedMemoryHuman = parsed["used_memory_human"]
	stats.ConnectedClients, _ = parseInt64(parsed["connected_clients"])
	stats.KeyspaceHits, _ = parseInt64(parsed["keyspace_hits"])
	stats.KeyspaceMisses, _ = parseInt64(parsed["keyspace_misses"])
	if hitsPlusMisses := stats.KeyspaceHits + stats.KeyspaceMisses; hitsPlusMisses > 0 {
		stats.HitRate = float64(stats.KeyspaceHits) / float64(hitsPlusMisses)
	}

	if dbSize, err := s.redisClient.DBSize(ctx); err == nil {
		stats.DBSize = dbSize
	}

	return &stats, nil
}

// ListOnlineUsers lists users considered online based on recent activity.
func (s *Service) ListOnlineUsers(ctx context.Context) ([]*vo.OnlineUserVO, error) {
	if s.redisClient == nil {
		return []*vo.OnlineUserVO{}, nil
	}

	keys, err := s.redisClient.Keys(ctx, "online:user:*")
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return []*vo.OnlineUserVO{}, nil
	}

	result := make([]*vo.OnlineUserVO, 0, len(keys))
	for _, key := range keys {
		val, err := s.redisClient.Get(ctx, key)
		if err != nil || val == "" {
			continue
		}
		var item vo.OnlineUserVO
		if err := json.Unmarshal([]byte(val), &item); err != nil {
			continue
		}
		result = append(result, &item)
	}
	return result, nil
}

// UpdateOnlineUser updates online user heartbeat info.
func (s *Service) UpdateOnlineUser(ctx context.Context, userCode, username, ip, userAgent string) error {
	if s.redisClient == nil || userCode == "" {
		return nil
	}

	item := vo.OnlineUserVO{
		UserCode:   userCode,
		Username:   username,
		IP:         ip,
		UserAgent:  userAgent,
		LastActive: time.Now().Unix(),
	}
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	key := "online:user:" + userCode
	return s.redisClient.Set(ctx, key, b, s.onlineTTL)
}

// parseRedisInfo parses the output of the INFO command into a map.
func parseRedisInfo(info string) map[string]string {
	res := make(map[string]string)
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		res[parts[0]] = strings.TrimSpace(parts[1])
	}
	return res
}

func parseInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}
