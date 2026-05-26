package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

func registerOperationLogRoutes(rg *gin.RouterGroup, permissionSvc interfaces.PermissionService, dataScopeSvc interfaces.DataScopeService, logSvc interfaces.OperationLogService) {
	if logSvc == nil {
		return
	}
	// 操作日志分页查询
	rg.GET("/operation-logs", middleware.RequirePermission(permissionSvc, "system:operation-log:list"), func(c *gin.Context) {
		listOperationLogs(c, logSvc, dataScopeSvc)
	})
}

func listOperationLogs(c *gin.Context, logSvc interfaces.OperationLogService, dataScopeSvc interfaces.DataScopeService) {
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

	var scope *interfaces.UserDataScope
	if dataScopeSvc != nil {
		claimsVal, _ := c.Get(middleware.CurrentUserCtxKey)
		if claims, ok := claimsVal.(*jwtutil.Claims); ok && claims != nil {
			resolved, err := dataScopeSvc.ResolveUserScope(c, claims.UserCode, claims.IsSuperAdmin)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			scope = resolved
		}
	}
	result, err := logSvc.PageWithScope(c, &req, scope)
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
