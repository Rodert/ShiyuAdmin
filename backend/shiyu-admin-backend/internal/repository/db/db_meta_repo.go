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
	TableName string `gorm:"column:table_name"`
	TableType string `gorm:"column:table_type"`
}

type columnRow struct {
	ColumnName    string  `gorm:"column:column_name"`
	DataType      string  `gorm:"column:data_type"`
	IsNullable    string  `gorm:"column:is_nullable"`
	MaxLength     *int64  `gorm:"column:character_maximum_length"`
	ColumnDefault *string `gorm:"column:column_default"`
}

func (r *DBMetaRepository) ListTables(ctx context.Context) ([]*interfaces.TableMeta, error) {
	rows := make([]*tableRow, 0)
	query := `SELECT table_name, table_type FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name`
	if err := r.db.WithContext(ctx).Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*interfaces.TableMeta, 0, len(rows))
	for _, row := range rows {
		result = append(result, &interfaces.TableMeta{
			TableName: row.TableName,
			TableType: row.TableType,
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
	query := `SELECT column_name, data_type, is_nullable, character_maximum_length, column_default FROM information_schema.columns WHERE table_schema = 'public' AND table_name = ? ORDER BY ordinal_position`
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
		})
	}
	return result, nil
}
