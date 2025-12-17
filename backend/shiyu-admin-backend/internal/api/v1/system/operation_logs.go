package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerOperationLogRoutes(rg *gin.RouterGroup, logSvc interfaces.OperationLogService) {
	if logSvc == nil {
		return
	}
	// 操作日志分页查询
	rg.GET("/operation-logs", func(c *gin.Context) {
		listOperationLogs(c, logSvc)
	})
}

func listOperationLogs(c *gin.Context, logSvc interfaces.OperationLogService) {
	var req dto.OperationLogPageRequest
	req.Page = 1
	req.PageSize = 10
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	result, err := logSvc.Page(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil {
		result = &vo.PageResult[vo.OperationLogVO]{
			Items: []*vo.OperationLogVO{},
			Page:  req.Page,
			Size:  req.PageSize,
			Total: 0,
		}
	}
	response.Success(c, result)
}
