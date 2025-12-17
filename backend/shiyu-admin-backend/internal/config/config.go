package config

import (
	"fmt"
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
	return &cfg, nil
}
