package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/model/vo"
	monitorinterfaces "shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/redis"
)

// Service implements MonitorService.
type Service struct {
	db          *gorm.DB
	redisClient *redis.Client
	// onlineTTL defines how long a user is considered online since last activity.
	onlineTTL  time.Duration
	revokedTTL time.Duration
}

// New creates a new MonitorService.
func New(redisClient *redis.Client, onlineTTL time.Duration, db *gorm.DB, revokedTTL ...time.Duration) monitorinterfaces.MonitorService {
	tokenTTL := onlineTTL
	if len(revokedTTL) > 0 && revokedTTL[0] > 0 {
		tokenTTL = revokedTTL[0]
	}
	return &Service{
		db:          db,
		redisClient: redisClient,
		onlineTTL:   onlineTTL,
		revokedTTL:  tokenTTL,
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

// GetDatabaseStats returns basic database health and schema statistics.
func (s *Service) GetDatabaseStats(ctx context.Context) (*vo.DatabaseStatsVO, error) {
	if s.db == nil {
		return &vo.DatabaseStatsVO{Status: "disabled"}, nil
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return nil, err
	}

	stats := &vo.DatabaseStatsVO{
		Status: "normal",
		Driver: s.db.Dialector.Name(),
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		stats.Status = "abnormal"
		return stats, nil
	}

	poolStats := sqlDB.Stats()
	stats.OpenConnections = poolStats.OpenConnections
	stats.InUse = poolStats.InUse
	stats.Idle = poolStats.Idle

	_ = s.db.WithContext(ctx).Raw(databaseNameSQL(stats.Driver)).Scan(&stats.Database).Error
	_ = s.db.WithContext(ctx).Raw(databaseVersionSQL(stats.Driver)).Scan(&stats.Version).Error
	_ = s.db.WithContext(ctx).Raw(tableCountSQL(stats.Driver)).Scan(&stats.TableCount).Error

	return stats, nil
}

// ListOnlineUsers lists users considered online based on recent activity.
func (s *Service) ListOnlineUsers(ctx context.Context) ([]*vo.OnlineUserVO, error) {
	if s.redisClient == nil {
		return []*vo.OnlineUserVO{}, nil
	}

	keys, err := s.redisClient.ScanKeys(ctx, "online:user:*", 200)
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
		if item.SessionID == "" {
			item.SessionID = strings.TrimPrefix(key, "online:user:")
		}
		result = append(result, &item)
	}
	if err := s.fillDepartmentNames(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateOnlineUser updates online user heartbeat info.
func (s *Service) UpdateOnlineUser(ctx context.Context, sessionID, userCode, username, ip, hostIP, userAgent string) error {
	if s.redisClient == nil || userCode == "" || sessionID == "" {
		return nil
	}

	now := time.Now().Unix()
	browser, browserVersion := parseBrowser(userAgent)
	item := vo.OnlineUserVO{
		SessionID:      sessionID,
		UserCode:       userCode,
		Username:       username,
		IP:             ip,
		LoginLocation:  parseLoginLocation(ip),
		HostIP:         hostIP,
		UserAgent:      userAgent,
		OS:             parseOS(userAgent),
		Browser:        browser,
		BrowserVersion: browserVersion,
		BrowserDetail:  browserDetail(browser, browserVersion),
		LoginTime:      now,
		LastActive:     now,
	}
	if existing, err := s.redisClient.Get(ctx, onlineUserKey(sessionID)); err == nil && existing != "" {
		var old vo.OnlineUserVO
		if err := json.Unmarshal([]byte(existing), &old); err == nil && old.LoginTime > 0 {
			item.LoginTime = old.LoginTime
		}
	}
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return s.redisClient.Set(ctx, onlineUserKey(sessionID), b, s.onlineTTL)
}

func (s *Service) ForceLogout(ctx context.Context, sessionID string) error {
	if s.redisClient == nil || sessionID == "" {
		return nil
	}
	if err := s.redisClient.Delete(ctx, onlineUserKey(sessionID)); err != nil {
		return err
	}
	return s.redisClient.Set(ctx, revokedSessionKey(sessionID), "1", s.revokedTTL)
}

func (s *Service) IsSessionRevoked(ctx context.Context, sessionID string) (bool, error) {
	if s.redisClient == nil || sessionID == "" {
		return false, nil
	}
	return s.redisClient.Exists(ctx, revokedSessionKey(sessionID))
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

func onlineUserKey(sessionID string) string {
	return "online:user:" + sessionID
}

func revokedSessionKey(sessionID string) string {
	return "online:revoked:" + sessionID
}

func (s *Service) fillDepartmentNames(ctx context.Context, users []*vo.OnlineUserVO) error {
	if s.db == nil || len(users) == 0 {
		return nil
	}
	codes := make([]string, 0, len(users))
	seen := make(map[string]struct{}, len(users))
	for _, user := range users {
		if user.UserCode == "" {
			continue
		}
		if _, ok := seen[user.UserCode]; ok {
			continue
		}
		seen[user.UserCode] = struct{}{}
		codes = append(codes, user.UserCode)
	}
	if len(codes) == 0 {
		return nil
	}

	rows := make([]struct {
		UserCode string `gorm:"column:user_code"`
		DeptName string `gorm:"column:dept_name"`
	}, 0, len(codes))
	if err := s.db.WithContext(ctx).
		Table("sys_users").
		Select("sys_users.user_code, sys_depts.dept_name").
		Joins("LEFT JOIN sys_depts ON sys_users.dept_code = sys_depts.dept_code").
		Where("sys_users.user_code IN ?", codes).
		Scan(&rows).Error; err != nil {
		return fmt.Errorf("fill online user departments failed: %w", err)
	}
	deptByUser := make(map[string]string, len(rows))
	for _, row := range rows {
		deptByUser[row.UserCode] = row.DeptName
	}
	for _, user := range users {
		user.DeptName = deptByUser[user.UserCode]
	}
	return nil
}

func parseOS(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		return "macOS"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "未知"
	}
}

func parseBrowser(userAgent string) (string, string) {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "edg/"):
		return "Microsoft Edge", userAgentVersion(userAgent, `(?i)Edg/([0-9.]+)`)
	case strings.Contains(ua, "crios/"):
		return "Chrome", userAgentVersion(userAgent, `(?i)CriOS/([0-9.]+)`)
	case strings.Contains(ua, "chrome/"):
		return "Chrome", userAgentVersion(userAgent, `(?i)Chrome/([0-9.]+)`)
	case strings.Contains(ua, "fxios/"):
		return "Firefox", userAgentVersion(userAgent, `(?i)FxiOS/([0-9.]+)`)
	case strings.Contains(ua, "firefox/"):
		return "Firefox", userAgentVersion(userAgent, `(?i)Firefox/([0-9.]+)`)
	case strings.Contains(ua, "safari/"):
		return "Safari", userAgentVersion(userAgent, `(?i)Version/([0-9.]+)`)
	default:
		return "未知", ""
	}
}

func userAgentVersion(userAgent string, pattern string) string {
	matches := regexp.MustCompile(pattern).FindStringSubmatch(userAgent)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func browserDetail(browser string, version string) string {
	if version == "" {
		return browser
	}
	return browser + " " + version
}

func parseLoginLocation(ipText string) string {
	ip := net.ParseIP(strings.TrimSpace(ipText))
	if ip == nil {
		return "未知地点"
	}
	if ip.IsLoopback() {
		return "本机"
	}
	if ip.IsPrivate() {
		return "内网"
	}
	return "公网IP（未配置归属地解析）"
}

func databaseNameSQL(driver string) string {
	switch driver {
	case "mysql":
		return "SELECT DATABASE()"
	case "sqlserver":
		return "SELECT DB_NAME()"
	case "sqlite":
		return "SELECT 'sqlite'"
	default:
		return "SELECT current_database()"
	}
}

func databaseVersionSQL(driver string) string {
	switch driver {
	case "mysql", "sqlite":
		return "SELECT VERSION()"
	case "sqlserver":
		return "SELECT @@VERSION"
	default:
		return "SHOW server_version"
	}
}

func tableCountSQL(driver string) string {
	switch driver {
	case "mysql":
		return "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()"
	case "sqlserver":
		return "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'dbo' AND table_type = 'BASE TABLE'"
	case "sqlite":
		return "SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'"
	default:
		return "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'"
	}
}
