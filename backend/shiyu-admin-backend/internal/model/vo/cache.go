package vo

// CacheDatabaseVO represents a Redis logical database summary.
type CacheDatabaseVO struct {
	DB       int   `json:"db"`
	KeyCount int64 `json:"key_count"`
}

// CacheKeyVO represents a Redis key summary.
type CacheKeyVO struct {
	Key string `json:"key"`
	// Type is one of string, list, set, zset, hash, stream, none.
	Type string `json:"type"`
	// TTL is seconds until expiration. -1 means no expiration, -2 means key does not exist.
	TTL int64 `json:"ttl"`
}

// CacheKeyListVO represents a key scan result.
type CacheKeyListVO struct {
	DB      int           `json:"db"`
	Pattern string        `json:"pattern"`
	Type    string        `json:"type"`
	Total   int           `json:"total"`
	Items   []*CacheKeyVO `json:"items"`
}

// CacheValueVO represents the value of a Redis key.
type CacheValueVO struct {
	DB    int         `json:"db"`
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	TTL   int64       `json:"ttl"`
	Value interface{} `json:"value"`
}
