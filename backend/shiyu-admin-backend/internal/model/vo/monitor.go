package vo

// CacheStatsVO represents basic Redis/cache statistics for monitoring.
type CacheStatsVO struct {
	RedisVersion      string  `json:"redis_version"`
	Mode              string  `json:"mode"`
	UsedMemory        int64   `json:"used_memory"`
	UsedMemoryHuman   string  `json:"used_memory_human"`
	DBSize            int64   `json:"db_size"`
	ConnectedClients  int64   `json:"connected_clients"`
	KeyspaceHits      int64   `json:"keyspace_hits"`
	KeyspaceMisses    int64   `json:"keyspace_misses"`
	HitRate           float64 `json:"hit_rate"`
}

// OnlineUserVO represents an online user session for monitoring.
type OnlineUserVO struct {
	UserCode   string `json:"user_code"`
	Username   string `json:"username"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	LastActive int64  `json:"last_active"`
}
