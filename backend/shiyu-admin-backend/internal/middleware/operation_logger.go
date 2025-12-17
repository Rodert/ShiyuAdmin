package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
)

// OperationLogger records write operations (POST/PUT/PATCH/DELETE) to operation logs.
func OperationLogger(logSvc interfaces.OperationLogService) gin.HandlerFunc {
	if logSvc == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		method := c.Request.Method
		if !isWriteMethod(method) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		module, action := deriveModuleAndAction(path, method)

		var userCode, username string
		if claimsVal, exists := c.Get(CurrentUserCtxKey); exists {
			if claims, ok := claimsVal.(*jwtutil.Claims); ok {
				userCode = claims.UserCode
				username = claims.Username
			}
		}

		statusCode := c.Writer.Status()
		status := 1
		var errorMsg string
		if statusCode >= http.StatusBadRequest {
			status = 0
			if len(c.Errors) > 0 {
				// Gin 会把 error 累积在 c.Errors 中
				msg := c.Errors.String()
				if len(msg) > 500 {
					msg = msg[:500]
				}
				errorMsg = msg
			}
		}

		logEntry := &entity.OperationLog{
			UserCode:  userCode,
			Username:  username,
			Module:    module,
			Action:    action,
			Method:    method,
			Path:      path,
			IP:        c.ClientIP(),
			Status:    status,
			ErrorMsg:  errorMsg,
			LatencyMs: latency,
		}

		ctx := context.WithoutCancel(c.Request.Context())
		go func(ctx context.Context, entry *entity.OperationLog) {
			_ = logSvc.Create(ctx, entry)
		}(ctx, logEntry)
	}
}

func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func deriveModuleAndAction(path, method string) (string, string) {
	module := ""
	const prefix = "/api/v1/system/"
	if strings.HasPrefix(path, prefix) {
		rest := strings.TrimPrefix(path, prefix)
		parts := strings.Split(rest, "/")
		if len(parts) > 0 && parts[0] != "" {
			module = "system-" + parts[0]
		}
	}

	action := "other"
	switch method {
	case http.MethodPost:
		action = "create"
	case http.MethodPut, http.MethodPatch:
		action = "update"
	case http.MethodDelete:
		action = "delete"
	}
	return module, action
}
