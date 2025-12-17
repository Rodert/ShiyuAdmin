package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerRoleDeptRoutes(rg *gin.RouterGroup, roleDeptSvc interfaces.RoleDeptService) {
	if roleDeptSvc == nil {
		return
	}
	rg.GET("/roles/:code/depts", func(c *gin.Context) {
		getRoleDepts(c, roleDeptSvc)
	})
	rg.PUT("/roles/:code/depts", func(c *gin.Context) {
		setRoleDepts(c, roleDeptSvc)
	})
}

func getRoleDepts(c *gin.Context, svc interfaces.RoleDeptService) {
	depts, err := svc.GetRoleDepts(c, c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.DeptVO, 0, len(depts))
	for _, d := range depts {
		items = append(items, vo.BuildDeptVO(d))
	}
	response.Success(c, items)
}

func setRoleDepts(c *gin.Context, svc interfaces.RoleDeptService) {
	var req dto.SetRoleDeptsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := svc.SetRoleDepts(c, c.Param("code"), req.DeptCodes); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}

