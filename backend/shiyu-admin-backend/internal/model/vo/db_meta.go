package vo

type TableMetaVO struct {
	TableName string `json:"table_name"`
	TableType string `json:"table_type"`
}

type ColumnMetaVO struct {
	ColumnName    string  `json:"column_name"`
	DataType      string  `json:"data_type"`
	IsNullable    bool    `json:"is_nullable"`
	MaxLength     *int64  `json:"max_length,omitempty"`
	ColumnDefault *string `json:"column_default,omitempty"`
}

// TableDataPageVO represents paginated table data rows.
type TableDataPageVO struct {
	Items []map[string]interface{} `json:"items"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
	Total int64                   `json:"total"`
}
