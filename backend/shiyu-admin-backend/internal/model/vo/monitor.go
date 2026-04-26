package vo

// CacheStatsVO represents basic Redis/cache statistics for monitoring.
type CacheStatsVO struct {
	RedisVersion     string  `json:"redis_version"`
	Mode             string  `json:"mode"`
	UsedMemory       int64   `json:"used_memory"`
	UsedMemoryHuman  string  `json:"used_memory_human"`
	DBSize           int64   `json:"db_size"`
	ConnectedClients int64   `json:"connected_clients"`
	KeyspaceHits     int64   `json:"keyspace_hits"`
	KeyspaceMisses   int64   `json:"keyspace_misses"`
	HitRate          float64 `json:"hit_rate"`
}

// OnlineUserVO represents an online user session for monitoring.
type OnlineUserVO struct {
	UserCode   string `json:"user_code"`
	Username   string `json:"username"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	LastActive int64  `json:"last_active"`
}

// DatabaseStatsVO represents database health and basic schema statistics.
type DatabaseStatsVO struct {
	Status          string `json:"status"`
	Driver          string `json:"driver"`
	Database        string `json:"database"`
	Version         string `json:"version"`
	TableCount      int64  `json:"table_count"`
	OpenConnections int    `json:"open_connections"`
	InUse           int    `json:"in_use"`
	Idle            int    `json:"idle"`
}
