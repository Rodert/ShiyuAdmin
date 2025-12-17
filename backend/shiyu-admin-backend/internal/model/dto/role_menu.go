package dto

// SetRoleMenusRequest defines role-menu assignment payload.
type SetRoleMenusRequest struct {
	MenuCodes []string `json:"menu_codes" binding:"required"`
}

