package interfaces

import (
	"context"

	"shiyu-admin-backend/internal/model/vo"
)

type DataManageService interface {
	VerifyCredential(ctx context.Context, username string, password string) bool
	ListTables(ctx context.Context) ([]*vo.TableMetaVO, error)
	ListColumns(ctx context.Context, tableName string) ([]*vo.ColumnMetaVO, error)
	PageTableData(ctx context.Context, tableName string, page, pageSize int) (*vo.TableDataPageVO, error)
}
