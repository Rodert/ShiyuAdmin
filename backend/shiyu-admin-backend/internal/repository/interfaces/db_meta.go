package interfaces

import "context"

type TableMeta struct {
	TableName string
	TableType string
}

type ColumnMeta struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	MaxLength     *int64
	ColumnDefault *string
}

type DBMetaRepository interface {
	ListTables(ctx context.Context) ([]*TableMeta, error)
	ListColumns(ctx context.Context, tableName string) ([]*ColumnMeta, error)
	ListRows(ctx context.Context, tableName string, page, pageSize int) ([]map[string]interface{}, int64, error)
}
