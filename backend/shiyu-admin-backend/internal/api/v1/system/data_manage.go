package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

// registerDataManageRoutes registers data management (DB metadata) routes.
func registerDataManageRoutes(rg *gin.RouterGroup, svc interfaces.DataManageService) {
	if svc == nil {
		return
	}

	// 列出所有表
	rg.GET("/data/tables", func(c *gin.Context) {
		listTables(c, svc)
	})

	// 列出指定表的字段
	rg.GET("/data/tables/:table/columns", func(c *gin.Context) {
		listTableColumns(c, svc)
	})

	// 分页查询指定表的数据
	rg.GET("/data/tables/:table/rows", func(c *gin.Context) {
		listTableRows(c, svc)
	})
}

func listTables(c *gin.Context, svc interfaces.DataManageService) {
	items, err := svc.ListTables(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*vo.TableMetaVO{}
	}
	response.Success(c, items)
}

// listTableRows returns paginated rows for a given table.
func listTableRows(c *gin.Context, svc interfaces.DataManageService) {
	table := c.Param("table")
	if table == "" {
		response.Error(c, http.StatusBadRequest, "表名不能为空")
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("page_size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		size = 10
	}

	result, err := svc.PageTableData(c, table, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		result = &vo.TableDataPageVO{
			Items: []map[string]interface{}{},
			Page:  page,
			Size:  size,
			Total: 0,
		}
	}
	response.Success(c, result)
}

func listTableColumns(c *gin.Context, svc interfaces.DataManageService) {
	table := c.Param("table")
	if table == "" {
		response.Error(c, http.StatusBadRequest, "表名不能为空")
		return
	}
	items, err := svc.ListColumns(c, table)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*vo.ColumnMetaVO{}
	}
	response.Success(c, items)
}
