package dto

// CreateMenuRequest defines menu creation payload.
type CreateMenuRequest struct {
	MenuCode   string `json:"menu_code" binding:"required"`
	ParentCode string `json:"parent_code"`
	MenuType   string `json:"menu_type" binding:"required"` // M=目录, C=菜单, F=按钮
	MenuName   string `json:"menu_name"`
	Perms      string `json:"perms"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Status     int    `json:"status"`
}

// UpdateMenuRequest defines menu update payload.
type UpdateMenuRequest struct {
	ParentCode *string `json:"parent_code"`
	MenuType   *string `json:"menu_type"`
	MenuName   *string `json:"menu_name"`
	Perms      *string `json:"perms"`
	Path       *string `json:"path"`
	Component  *string `json:"component"`
	Status     *int    `json:"status"`
}

