package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ForUpdate locks selected rows until the surrounding transaction ends.
// SQL Server requires table hints; SQLite uses its transaction-level writer lock.
func ForUpdate(db *gorm.DB, table string) *gorm.DB {
	if db.Dialector.Name() == "sqlserver" {
		return db.Table(table).Table("? WITH (UPDLOCK, ROWLOCK)", clause.Table{Name: table})
	}
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}
