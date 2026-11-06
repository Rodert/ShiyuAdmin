package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadOrCreateSecret creates a private, persistent secret and never overwrites it.
func LoadOrCreateSecret(path string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return "", err
	}
	if b, err := os.ReadFile(path); err == nil {
		value := strings.TrimSpace(string(b))
		if len(value) < 32 {
			return "", fmt.Errorf("invalid secret file: %s", path)
		}
		return value, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	value := hex.EncodeToString(b)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("create secret file (retry startup if another instance is initializing): %w", err)
	}
	_, writeErr := f.WriteString(value + "\n")
	if writeErr == nil {
		writeErr = f.Sync()
	}
	closeErr := f.Close()
	if writeErr != nil {
		return "", writeErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	return value, nil
}

func PrepareSecurity(cfg *Config) error {
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}
	if err := os.MkdirAll(filepath.Join(cfg.DataDir, "uploads"), 0750); err != nil {
		return err
	}
	if cfg.JWT.ExpireTime <= 0 {
		return fmt.Errorf("jwt.expire_time must be positive")
	}
	if cfg.JWT.Secret == "" {
		secret, err := LoadOrCreateSecret(filepath.Join(cfg.DataDir, "jwt.key"))
		if err != nil {
			return err
		}
		cfg.JWT.Secret = secret
	} else if len(cfg.JWT.Secret) < 32 || strings.Contains(cfg.JWT.Secret, "shiyu-admin-scaffold") || strings.Contains(cfg.JWT.Secret, "change-me") {
		return fmt.Errorf("unsafe JWT secret: remove the default to generate a persistent key, or set SHIYU_JWT_SECRET to a random secret of at least 32 bytes")
	}
	if cfg.Auth.UseMock && cfg.Server.Mode == "release" {
		return fmt.Errorf("mock authentication is forbidden in release mode")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 12 || len(password) > 72 {
		return fmt.Errorf("密码须为 12–72 字节")
	}
	if password == "Admin@123456" || password == "Password123!" || password == "123456789012" {
		return fmt.Errorf("请使用独立的强密码")
	}
	return nil
}
