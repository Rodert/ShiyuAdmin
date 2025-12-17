package data_manage

import (
	"context"
	"fmt"
	"strings"

	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/repository/interfaces"
	serviceinterfaces "shiyu-admin-backend/internal/service/interfaces"
)

type Service struct {
	repo interfaces.DBMetaRepository
}

func New(repo interfaces.DBMetaRepository) serviceinterfaces.DataManageService {
	return &Service{repo: repo}
}

func (s *Service) ListTables(ctx context.Context) ([]*vo.TableMetaVO, error) {
	tables, err := s.repo.ListTables(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.TableMetaVO, 0, len(tables))
	for _, t := range tables {
		result = append(result, &vo.TableMetaVO{
			TableName: t.TableName,
			TableType: t.TableType,
		})
	}
	return result, nil
}

// PageTableData returns paginated data rows for a table.
func (s *Service) PageTableData(ctx context.Context, tableName string, page, pageSize int) (*vo.TableDataPageVO, error) {
	if tableName == "" {
		return nil, fmt.Errorf("table name is required")
	}
	// Basic validation: ensure table exists in metadata list.
	tables, err := s.repo.ListTables(ctx)
	if err != nil {
		return nil, err
	}
	found := false
	for _, t := range tables {
		if t.TableName == tableName {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("table %s not found", tableName)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	rows, total, err := s.repo.ListRows(ctx, tableName, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &vo.TableDataPageVO{
		Items: rows,
		Page:  page,
		Size:  pageSize,
		Total: total,
	}, nil
}

func (s *Service) ListColumns(ctx context.Context, tableName string) ([]*vo.ColumnMetaVO, error) {
	cols, err := s.repo.ListColumns(ctx, tableName)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.ColumnMetaVO, 0, len(cols))
	for _, c := range cols {
		nullable := strings.EqualFold(c.IsNullable, "YES")
		result = append(result, &vo.ColumnMetaVO{
			ColumnName:    c.ColumnName,
			DataType:      c.DataType,
			IsNullable:    nullable,
			MaxLength:     c.MaxLength,
			ColumnDefault: c.ColumnDefault,
		})
	}
	return result, nil
}
