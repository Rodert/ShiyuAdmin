package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/pkg/logger"
)

const traceIDHeader = "X-Trace-Id"

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(traceIDHeader)
		if traceID == "" {
			traceID = newTraceID()
		}
		ctx := logger.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Header(traceIDHeader, traceID)
		c.Set("trace_id", traceID)
		c.Next()
	}
}

// RequestLogger records metadata only. Bodies may contain passwords, tokens or
// multi-gigabyte files and must never be buffered by logging middleware.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		logger.Info(c.Request.Context(), "http_request",
			"method", c.Request.Method, "path", path, "status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(), "ip", c.ClientIP())
	}
}

func newTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
