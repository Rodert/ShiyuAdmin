package database

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"shiyu-admin-backend/internal/model/entity"
)

// This checks generated syntax without pretending to replace the live SQL Server CI suite.
func TestSQLServerUpdateLockSyntax(t *testing.T) {
	db, err := gorm.Open(sqlserver.Open("sqlserver://unused:unused@127.0.0.1:1?database=unused"), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	require.NoError(t, err)
	stmt := ForUpdate(db, "sys_media_files").Where("file_code = ?", "file").First(&entity.MediaFile{}).Statement
	require.NoError(t, stmt.Error)
	require.Contains(t, stmt.SQL.String(), "WITH (UPDLOCK, ROWLOCK)")
	require.NotContains(t, stmt.SQL.String(), "FOR UPDATE")
	require.Contains(t, stmt.SQL.String(), "sys_media_files")
	require.Equal(t, []interface{}{"file"}, stmt.Vars)
}
