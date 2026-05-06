package dto

// ListUserRequest captures pagination params.
type ListUserRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// CreateUserRequest defines user creation payload.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	DeptCode string `json:"dept_code"`
	Status   int    `json:"status"`
}

// UpdateUserRequest defines user update payload.
type UpdateUserRequest struct {
	Nickname *string `json:"nickname"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
	DeptCode *string `json:"dept_code"`
	Status   *int    `json:"status"`
	Password *string `json:"password"`
}

// UpdateProfileRequest defines editable fields for the current user.
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
}

// ChangePasswordRequest defines current-user password change payload.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Role DTOs
type ListRoleRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type CreateRoleRequest struct {
	RoleCode  string `json:"role_code" binding:"required"`
	RoleName  string `json:"role_name" binding:"required"`
	RoleKey   string `json:"role_key" binding:"required"`
	DataScope string `json:"data_scope"`
	Status    int    `json:"status"`
}

type UpdateRoleRequest struct {
	RoleName  *string `json:"role_name"`
	DataScope *string `json:"data_scope"`
	Status    *int    `json:"status"`
}
