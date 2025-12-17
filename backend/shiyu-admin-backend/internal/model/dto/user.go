package dto

// ListUserRequest captures pagination params.
type ListUserRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// CreateUserRequest defines user creation payload.
type CreateUserRequest struct {
	UserCode string `json:"user_code" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DeptCode string `json:"dept_code"`
	Status   int    `json:"status"`
}

// UpdateUserRequest defines user update payload.
type UpdateUserRequest struct {
	Nickname *string `json:"nickname"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	DeptCode *string `json:"dept_code"`
	Status   *int    `json:"status"`
	Password *string `json:"password"`
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
	RoleKey   *string `json:"role_key"`
	DataScope *string `json:"data_scope"`
	Status    *int    `json:"status"`
}
