package dto

// SetUserRolesRequest defines user-role assignment payload.
type SetUserRolesRequest struct {
	RoleCodes []string `json:"role_codes" binding:"required"`
}

