package db

import (
	"context"

	"gorm.io/gorm"

	"shiyu-admin-backend/internal/repository/interfaces"
)

type DBMetaRepository struct {
	db *gorm.DB
}

func NewDBMetaRepository(db *gorm.DB) interfaces.DBMetaRepository {
	return &DBMetaRepository{db: db}
}

type tableRow struct {
	TableName    string `gorm:"column:table_name"`
	TableType    string `gorm:"column:table_type"`
	TableComment string `gorm:"column:table_comment"`
}

type columnRow struct {
	ColumnName    string  `gorm:"column:column_name"`
	DataType      string  `gorm:"column:data_type"`
	IsNullable    string  `gorm:"column:is_nullable"`
	MaxLength     *int64  `gorm:"column:character_maximum_length"`
	ColumnDefault *string `gorm:"column:column_default"`
	ColumnComment string  `gorm:"column:column_comment"`
}

func (r *DBMetaRepository) ListTables(ctx context.Context) ([]*interfaces.TableMeta, error) {
	rows := make([]*tableRow, 0)
	query := r.listTablesSQL()
	if err := r.db.WithContext(ctx).Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*interfaces.TableMeta, 0, len(rows))
	for _, row := range rows {
		result = append(result, &interfaces.TableMeta{
			TableName:    row.TableName,
			TableType:    row.TableType,
			TableComment: row.TableComment,
		})
	}
	return result, nil
}

// ListRows returns paginated data rows for a given table.
func (r *DBMetaRepository) ListRows(ctx context.Context, tableName string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var total int64
	query := r.db.WithContext(ctx).Table(tableName)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []map[string]interface{}{}, 0, nil
	}
	offset := (page - 1) * pageSize
	rows := make([]map[string]interface{}, 0)
	if err := query.
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *DBMetaRepository) ListColumns(ctx context.Context, tableName string) ([]*interfaces.ColumnMeta, error) {
	rows := make([]*columnRow, 0)
	query := r.listColumnsSQL()
	if err := r.db.WithContext(ctx).Raw(query, tableName).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*interfaces.ColumnMeta, 0, len(rows))
	for _, row := range rows {
		result = append(result, &interfaces.ColumnMeta{
			ColumnName:    row.ColumnName,
			DataType:      row.DataType,
			IsNullable:    row.IsNullable,
			MaxLength:     row.MaxLength,
			ColumnDefault: row.ColumnDefault,
			ColumnComment: row.ColumnComment,
		})
	}
	return result, nil
}

func (r *DBMetaRepository) listTablesSQL() string {
	switch r.db.Dialector.Name() {
	case "mysql":
		return `
SELECT table_name, table_type, table_comment
FROM information_schema.tables
WHERE table_schema = DATABASE()
ORDER BY table_name`
	case "sqlite":
		return `
SELECT name AS table_name, type AS table_type, '' AS table_comment
FROM sqlite_master
WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%'
ORDER BY name`
	default:
		return `
SELECT
  t.table_name,
  t.table_type,
  COALESCE(obj_description(c.oid), '') AS table_comment
FROM information_schema.tables t
LEFT JOIN pg_class c ON c.relname = t.table_name
LEFT JOIN pg_namespace n ON n.oid = c.relnamespace AND n.nspname = t.table_schema
WHERE t.table_schema = 'public'
ORDER BY t.table_name`
	}
}

func (r *DBMetaRepository) listColumnsSQL() string {
	switch r.db.Dialector.Name() {
	case "mysql":
		return `
SELECT column_name, data_type, is_nullable, character_maximum_length, column_default, column_comment
FROM information_schema.columns
WHERE table_schema = DATABASE() AND table_name = ?
ORDER BY ordinal_position`
	case "sqlite":
		return `
SELECT name AS column_name, type AS data_type, "YES" AS is_nullable, NULL AS character_maximum_length, dflt_value AS column_default, '' AS column_comment
FROM pragma_table_info(?)
ORDER BY cid`
	default:
		return `
SELECT
  c.column_name,
  c.data_type,
  c.is_nullable,
  c.character_maximum_length,
  c.column_default,
  COALESCE(d.description, '') AS column_comment
FROM information_schema.columns c
LEFT JOIN pg_class pc ON pc.relname = c.table_name
LEFT JOIN pg_namespace pn ON pn.oid = pc.relnamespace AND pn.nspname = c.table_schema
LEFT JOIN pg_attribute pa ON pa.attrelid = pc.oid AND pa.attname = c.column_name
LEFT JOIN pg_description d ON d.objoid = pc.oid AND d.objsubid = pa.attnum
WHERE c.table_schema = 'public' AND c.table_name = ?
ORDER BY c.ordinal_position`
	}
}
