package database

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"shiyu-admin-backend/internal/config"
)

// Connect initializes a gorm DB based on config.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("nil config")
	}

	var dialector gorm.Dialector
	switch cfg.Database.Driver {
	case "postgres", "postgresql":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			cfg.Database.Host, cfg.Database.Username, cfg.Database.Password,
			cfg.Database.Database, cfg.Database.Port, cfg.Database.SSLMode, cfg.Database.Timezone,
		)
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.Username, cfg.Database.Password,
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Database,
		)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.Database)
	default:
		return nil, errors.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	gcfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}
	db, err := gorm.Open(dialector, gcfg)
	if err != nil {
		return nil, errors.Wrap(err, "open database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get sql.DB")
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}



