package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerDeptRoutes(rg *gin.RouterGroup, deptSvc interfaces.DeptService) {
	if deptSvc == nil {
		return
	}
	rg.GET("/depts", func(c *gin.Context) {
		listDepts(c, deptSvc)
	})
	rg.GET("/depts/tree", func(c *gin.Context) {
		listDeptTree(c, deptSvc)
	})
	rg.GET("/depts/:code", func(c *gin.Context) {
		getDept(c, deptSvc)
	})
	rg.POST("/depts", func(c *gin.Context) {
		createDept(c, deptSvc)
	})
	rg.PUT("/depts/:code", func(c *gin.Context) {
		updateDept(c, deptSvc)
	})
	rg.DELETE("/depts/:code", func(c *gin.Context) {
		deleteDept(c, deptSvc)
	})
}

func listDepts(c *gin.Context, svc interfaces.DeptService) {
	depts, err := svc.List(c)
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

func listDeptTree(c *gin.Context, svc interfaces.DeptService) {
	depts, err := svc.ListTree(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	tree := vo.BuildDeptTree(depts)
	response.Success(c, tree)
}

func getDept(c *gin.Context, svc interfaces.DeptService) {
	dept, err := svc.GetByCode(c, c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if dept == nil {
		response.Error(c, http.StatusNotFound, "部门不存在")
		return
	}
	response.Success(c, vo.BuildDeptVO(dept))
}

func createDept(c *gin.Context, svc interfaces.DeptService) {
	var req dto.CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	dept, err := svc.Create(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, vo.BuildDeptVO(dept))
}

func updateDept(c *gin.Context, svc interfaces.DeptService) {
	var req dto.UpdateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	dept, err := svc.Update(c, c.Param("code"), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if dept == nil {
		response.Error(c, http.StatusNotFound, "部门不存在")
		return
	}
	response.Success(c, vo.BuildDeptVO(dept))
}

func deleteDept(c *gin.Context, svc interfaces.DeptService) {
	if err := svc.Delete(c, c.Param("code")); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

