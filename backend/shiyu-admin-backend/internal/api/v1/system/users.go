package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/response"
)

func registerUserRoutes(rg *gin.RouterGroup, userSvc interfaces.UserService) {
	if userSvc == nil {
		return
	}
	rg.GET("/users", func(c *gin.Context) {
		listUsers(c, userSvc)
	})
	rg.POST("/users", func(c *gin.Context) {
		createUser(c, userSvc)
	})
	rg.PUT("/users/:code", func(c *gin.Context) {
		updateUser(c, userSvc)
	})
	rg.DELETE("/users/:code", func(c *gin.Context) {
		deleteUser(c, userSvc)
	})
}

func listUsers(c *gin.Context, userSvc interfaces.UserService) {
	var req dto.ListUserRequest
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

	users, total, err := userSvc.List(c, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.UserVO, 0, len(users))
	for _, u := range users {
		items = append(items, vo.BuildUserVO(u))
	}
	response.Success(c, vo.PageResult[vo.UserVO]{
		Items: items,
		Page:  req.Page,
		Size:  req.PageSize,
		Total: total,
	})
}

func createUser(c *gin.Context, userSvc interfaces.UserService) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := userSvc.Create(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, vo.BuildUserVO(user))
}

func updateUser(c *gin.Context, userSvc interfaces.UserService) {
	code := c.Param("code")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	user, err := userSvc.Update(c, code, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	response.Success(c, vo.BuildUserVO(user))
}

func deleteUser(c *gin.Context, userSvc interfaces.UserService) {
	code := c.Param("code")
	if err := userSvc.Delete(c, code); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
