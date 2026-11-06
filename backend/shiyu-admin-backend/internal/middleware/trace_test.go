package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"os"
	"path/filepath"
	"shiyu-admin-backend/pkg/logger"
	"strings"
	"testing"
)

func TestRequestLoggerDoesNotRecordCredentials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "requests.log")
	logger.Init("info", "json", 1, path, 1)
	t.Cleanup(func() { logger.Init("info", "json", 0, "", 0) })
	r := gin.New()
	r.Use(RequestLogger())
	r.POST("/login", func(c *gin.Context) { c.JSON(200, gin.H{"token": "sensitive-response-token"}) })
	req := httptest.NewRequest("POST", "/login?token=sensitive-query-token", strings.NewReader(`{"password":"sensitive-password"}`))
	req.Header.Set("Authorization", "Bearer sensitive-header-token")
	r.ServeHTTP(httptest.NewRecorder(), req)
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(b), "http_request")
	require.NotContains(t, string(b), "sensitive-")
}
