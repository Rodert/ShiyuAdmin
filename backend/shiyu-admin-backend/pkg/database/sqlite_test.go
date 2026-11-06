package database

import (
	"path/filepath"
	"testing"

	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/config"
)

func TestSQLiteTransactionReservesWriterBeforeRead(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = filepath.Join(t.TempDir(), "writers.db")
	connect := func() *gorm.DB {
		db, err := Connect(cfg)
		require.NoError(t, err)
		pool, err := db.DB()
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, pool.Close()) })
		return db
	}
	first, second := connect(), connect()
	require.NoError(t, first.Exec("CREATE TABLE counters (id integer PRIMARY KEY, value integer NOT NULL)").Error)
	require.NoError(t, first.Exec("INSERT INTO counters (id, value) VALUES (1, 0)").Error)
	// Keep the contention assertion fast; normal connections retain the 5-second timeout.
	require.NoError(t, second.Exec("PRAGMA busy_timeout = 25").Error)
	tx := first.Begin()
	require.NoError(t, tx.Error)
	defer tx.Rollback()
	var value int
	require.NoError(t, tx.Raw("SELECT value FROM counters WHERE id = 1").Scan(&value).Error)
	require.Zero(t, value)
	// A deferred transaction lets this competing write commit, then cannot upgrade
	// its stale WAL snapshot to a writer. Immediate transactions reserve the writer first.
	err := second.Exec("UPDATE counters SET value = value + 1 WHERE id = 1").Error
	var busy sqlite3.Error
	require.ErrorAs(t, err, &busy)
	require.Equal(t, sqlite3.ErrBusy, busy.Code)
	require.NoError(t, tx.Exec("UPDATE counters SET value = ? WHERE id = 1", value+1).Error)
	require.NoError(t, tx.Commit().Error)
	// Once the first writer commits, another read/modify/write transaction sees its result.
	require.NoError(t, second.Transaction(func(next *gorm.DB) error {
		if err := next.Raw("SELECT value FROM counters WHERE id = 1").Scan(&value).Error; err != nil {
			return err
		}
		return next.Exec("UPDATE counters SET value = ? WHERE id = 1", value+1).Error
	}))
	require.NoError(t, first.Raw("SELECT value FROM counters WHERE id = 1").Scan(&value).Error)
	require.Equal(t, 2, value)
}
