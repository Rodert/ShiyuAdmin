package middleware

import (
	"github.com/gin-gonic/gin"
	"shiyu-admin-backend/pkg/response"
	"sync"
	"time"
)

// LoginLimit bounds password-hashing work per IP and bounds its own memory use.
func LoginLimit() gin.HandlerFunc {
	type bucket struct {
		start time.Time
		count int
	}
	var mu sync.Mutex
	buckets := map[string]bucket{}
	return func(c *gin.Context) {
		now, key := time.Now(), c.ClientIP()
		mu.Lock()
		for ip, entry := range buckets {
			if now.Sub(entry.start) >= time.Minute {
				delete(buckets, ip)
			}
		}
		entry, exists := buckets[key]
		blocked := entry.count >= 20 || (!exists && len(buckets) >= 4096)
		if !blocked {
			if !exists {
				entry.start = now
			}
			entry.count++
			buckets[key] = entry
		}
		mu.Unlock()
		if blocked {
			c.Header("Retry-After", "60")
			response.Error(c, 429, "登录尝试过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
