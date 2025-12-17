package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
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

func RequestLogger() gin.HandlerFunc {
	const maxBodyLogLength = 2048
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body (for non-GET) with limit; restore body so handlers can read it.
		var requestBody string
		method := c.Request.Method
		if method == "GET" || method == "HEAD" {
			requestBody = c.Request.URL.RawQuery
		} else if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				if len(bodyBytes) > maxBodyLogLength {
					requestBody = string(bodyBytes[:maxBodyLogLength]) + "..." 
				} else {
					requestBody = string(bodyBytes)
				}
			}
		}

		// Wrap ResponseWriter to capture response body
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start).Milliseconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()
		responseBody := blw.body.String()
		if len(responseBody) > maxBodyLogLength {
			responseBody = responseBody[:maxBodyLogLength] + "..." 
		}
		ctx := c.Request.Context()
		logger.Info(ctx, "http_request",
			"method", method,
			"path", path,
			"status", status,
			"latency_ms", latency,
			"ip", clientIP,
			"user_agent", ua,
			"request_body", requestBody,
			"response_body", responseBody,
		)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func newTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
