package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"shiyu-admin-backend/internal/api/v1/system"
	"shiyu-admin-backend/internal/config"
	"shiyu-admin-backend/internal/middleware"
	repoMock "shiyu-admin-backend/internal/repository/mock"
	authsvc "shiyu-admin-backend/internal/service/auth"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup mock services
	authRepo := repoMock.NewAuthMockRepository()
	cfg := &config.Config{
		JWT: config.Config{
			JWT: struct {
				Secret     string `mapstructure:"secret"`
				ExpireTime int64  `mapstructure:"expire_time"`
				Issuer     string `mapstructure:"issuer"`
			}{
				Secret:     "test-secret",
				ExpireTime: 3600,
				Issuer:     "test",
			},
		}.JWT,
	}
	authSvc := authsvc.New(authRepo, cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime)
	authMiddleware := middleware.Auth(cfg.JWT.Secret)

	api := router.Group("/api/v1")
	system.RegisterRoutes(api, authSvc, authMiddleware, nil, nil, nil, nil, nil, nil, nil)

	return router
}

func TestPing(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/system/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
}

func TestLogin(t *testing.T) {
	router := setupRouter()

	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	jsonData, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/v1/system/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := setupRouter()

	loginData := map[string]string{
		"username": "admin",
		"password": "wrongpassword",
	}
	jsonData, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/v1/system/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestProfileWithoutAuth(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/system/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestProfileWithAuth(t *testing.T) {
	router := setupRouter()

	// First login to get token
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}
	jsonData, _ := json.Marshal(loginData)

	loginReq, _ := http.NewRequest("POST", "/api/v1/system/auth/login", bytes.NewBuffer(jsonData))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	var loginResponse map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	// Use token to access profile
	req, _ := http.NewRequest("GET", "/api/v1/system/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
}

