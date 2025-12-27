package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	Server struct {
		Port         string        `mapstructure:"port"`
		Mode         string        `mapstructure:"mode"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	} `mapstructure:"server"`
	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		SSLMode  string `mapstructure:"ssl_mode"`
		Timezone string `mapstructure:"timezone"`
	} `mapstructure:"database"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	Auth struct {
		UseMock bool `mapstructure:"use_mock"`
	} `mapstructure:"auth"`
	Bootstrap struct {
		AdminUsername string `mapstructure:"admin_username"`
		AdminPassword string `mapstructure:"admin_password"`
		AdminNickname string `mapstructure:"admin_nickname"`
	} `mapstructure:"bootstrap"`
	JWT struct {
		Secret     string `mapstructure:"secret"`
		ExpireTime int64  `mapstructure:"expire_time"`
		Issuer     string `mapstructure:"issuer"`
	} `mapstructure:"jwt"`
	Log struct {
		Level         string `mapstructure:"level"`
		Format        string `mapstructure:"format"`
		RetentionDays int    `mapstructure:"retention_days"`
		FilePath      string `mapstructure:"file_path"`
		MaxSizeMB     int    `mapstructure:"max_size_mb"`
	} `mapstructure:"log"`
}

// Load loads configuration from file and environment variables.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// 支持 DATABASE_URL 环境变量（通用格式，适用于 Render、Railway 等平台）
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if err := parseDatabaseURL(dbURL, &cfg); err != nil {
			return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
		}
	}

	// 支持 REDIS_URL 环境变量
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if err := parseRedisURL(redisURL, &cfg); err != nil {
			return nil, fmt.Errorf("parse REDIS_URL: %w", err)
		}
	}

	// 支持 PORT 环境变量（Render 等平台会设置）
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	return &cfg, nil
}

// parseDatabaseURL 解析 DATABASE_URL 格式: postgresql://user:password@host:port/dbname?sslmode=require
func parseDatabaseURL(dbURL string, cfg *Config) error {
	u, err := url.Parse(dbURL)
	if err != nil {
		return fmt.Errorf("invalid database URL: %w", err)
	}

	cfg.Database.Driver = "postgres"
	cfg.Database.Host = u.Hostname()
	if u.Port() != "" {
		port, err := strconv.Atoi(u.Port())
		if err != nil {
			return fmt.Errorf("invalid port in database URL: %w", err)
		}
		cfg.Database.Port = port
	}
	cfg.Database.Username = u.User.Username()
	if password, ok := u.User.Password(); ok {
		cfg.Database.Password = password
	}
	cfg.Database.Database = strings.TrimPrefix(u.Path, "/")

	// 解析查询参数
	query := u.Query()
	if sslMode := query.Get("sslmode"); sslMode != "" {
		cfg.Database.SSLMode = sslMode
	} else {
		// 默认使用 require（生产环境安全）
		cfg.Database.SSLMode = "require"
	}

	return nil
}

// parseRedisURL 解析 REDIS_URL 格式: redis://:password@host:port/db 或 redis://host:port/db
func parseRedisURL(redisURL string, cfg *Config) error {
	u, err := url.Parse(redisURL)
	if err != nil {
		return fmt.Errorf("invalid redis URL: %w", err)
	}

	cfg.Redis.Host = u.Hostname()
	if u.Port() != "" {
		port, err := strconv.Atoi(u.Port())
		if err != nil {
			return fmt.Errorf("invalid port in redis URL: %w", err)
		}
		cfg.Redis.Port = port
	}
	if password, ok := u.User.Password(); ok {
		cfg.Redis.Password = password
	}
	// 解析数据库编号
	if dbStr := strings.TrimPrefix(u.Path, "/"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			cfg.Redis.DB = db
		}
	}

	return nil
}
