package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/middleware"
	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/service/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
	"shiyu-admin-backend/pkg/response"
)

func registerMenuRoutes(rg *gin.RouterGroup, menuSvc interfaces.MenuService, userRoleSvc interfaces.UserRoleService, roleMenuSvc interfaces.RoleMenuService) {
	if menuSvc == nil {
		return
	}
	rg.GET("/menus", func(c *gin.Context) {
		listMenus(c, menuSvc)
	})
	rg.GET("/menus/tree", func(c *gin.Context) {
		listMenuTree(c, menuSvc, userRoleSvc, roleMenuSvc)
	})
	rg.GET("/menus/:code", func(c *gin.Context) {
		getMenu(c, menuSvc)
	})
	rg.POST("/menus", func(c *gin.Context) {
		createMenu(c, menuSvc)
	})
	rg.PUT("/menus/:code", func(c *gin.Context) {
		updateMenu(c, menuSvc)
	})
	rg.DELETE("/menus/:code", func(c *gin.Context) {
		deleteMenu(c, menuSvc)
	})
}

func listMenus(c *gin.Context, svc interfaces.MenuService) {
	menus, err := svc.List(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	items := make([]*vo.MenuVO, 0, len(menus))
	for _, m := range menus {
		items = append(items, vo.BuildMenuVO(m))
	}
	response.Success(c, items)
}

func listMenuTree(c *gin.Context, menuSvc interfaces.MenuService, userRoleSvc interfaces.UserRoleService, roleMenuSvc interfaces.RoleMenuService) {
	menus, err := menuSvc.ListTree(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	filtered := menus
	if userRoleSvc != nil && roleMenuSvc != nil {
		if claimsVal, exists := c.Get(middleware.CurrentUserCtxKey); exists {
			if claims, ok := claimsVal.(*jwtutil.Claims); ok && claims.UserCode != "" {
				if claims.IsSuperAdmin {
					filtered = menus
				} else {
					allowedCodes, err := collectUserMenuCodes(c, userRoleSvc, roleMenuSvc, claims.UserCode)
					if err != nil {
						response.Error(c, http.StatusInternalServerError, err.Error())
						return
					}
					if len(allowedCodes) == 0 {
						filtered = []*entity.Menu{}
					} else {
						menuIndex := buildMenuIndex(menus)
						includeParentMenus(allowedCodes, menuIndex)
						filtered = filterMenusByCodes(menus, allowedCodes)
					}
				}
			}
		}
	}

	tree := vo.BuildMenuTree(filtered)
	response.Success(c, tree)
}

func getMenu(c *gin.Context, svc interfaces.MenuService) {
	menu, err := svc.GetByCode(c, c.Param("code"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if menu == nil {
		response.Error(c, http.StatusNotFound, "菜单不存在")
		return
	}
	response.Success(c, vo.BuildMenuVO(menu))
}

func createMenu(c *gin.Context, svc interfaces.MenuService) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	menu, err := svc.Create(c, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, vo.BuildMenuVO(menu))
}

func updateMenu(c *gin.Context, svc interfaces.MenuService) {
	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	menu, err := svc.Update(c, c.Param("code"), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if menu == nil {
		response.Error(c, http.StatusNotFound, "菜单不存在")
		return
	}
	response.Success(c, vo.BuildMenuVO(menu))
}

func deleteMenu(c *gin.Context, svc interfaces.MenuService) {
	if err := svc.Delete(c, c.Param("code")); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func collectUserMenuCodes(c *gin.Context, userRoleSvc interfaces.UserRoleService, roleMenuSvc interfaces.RoleMenuService, userCode string) (map[string]struct{}, error) {
	allowed := make(map[string]struct{})

	roles, err := userRoleSvc.GetUserRoles(c, userCode)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		menus, err := roleMenuSvc.GetRoleMenus(c, role.RoleCode)
		if err != nil {
			return nil, err
		}
		for _, menu := range menus {
			if menu == nil {
				continue
			}
			allowed[menu.MenuCode] = struct{}{}
		}
	}

	return allowed, nil
}

func buildMenuIndex(menus []*entity.Menu) map[string]*entity.Menu {
	index := make(map[string]*entity.Menu, len(menus))
	for _, m := range menus {
		if m == nil {
			continue
		}
		index[m.MenuCode] = m
	}
	return index
}

func includeParentMenus(allowed map[string]struct{}, menuIndex map[string]*entity.Menu) {
	for code := range allowed {
		includeParentChain(code, allowed, menuIndex)
	}
}

func includeParentChain(code string, allowed map[string]struct{}, menuIndex map[string]*entity.Menu) {
	current := menuIndex[code]
	if current == nil {
		return
	}
	parentCode := current.ParentCode
	for parentCode != "" {
		if _, exists := allowed[parentCode]; exists {
			break
		}
		parent := menuIndex[parentCode]
		if parent == nil {
			break
		}
		allowed[parentCode] = struct{}{}
		parentCode = parent.ParentCode
	}
}

func filterMenusByCodes(menus []*entity.Menu, allowed map[string]struct{}) []*entity.Menu {
	filtered := make([]*entity.Menu, 0, len(allowed))
	for _, menu := range menus {
		if menu == nil {
			continue
		}
		if _, ok := allowed[menu.MenuCode]; ok {
			filtered = append(filtered, menu)
		}
	}
	return filtered
}
