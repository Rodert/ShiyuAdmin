package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a permissive CORS middleware suitable for APIs; adjust as needed.
func CORS() gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AllowCredentials = true
	cfg.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Accept"}
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	cfg.MaxAge = 12 * time.Hour
	return cors.New(cfg)
}

