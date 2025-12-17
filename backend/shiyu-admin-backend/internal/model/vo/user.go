package vo

import "shiyu-admin-backend/internal/model/entity"

// UserVO basic user info for API responses.
type UserVO struct {
	UserCode string `json:"user_code"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DeptCode string `json:"dept_code"`
	Status   int    `json:"status"`
}

// PageResult wraps paginated data.
type PageResult[T any] struct {
	Items []*T  `json:"items"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

func BuildUserVO(u *entity.User) *UserVO {
	if u == nil {
		return nil
	}
	return &UserVO{
		UserCode: u.UserCode,
		Username: u.Username,
		Nickname: u.Nickname,
		Email:    u.Email,
		Phone:    u.Phone,
		DeptCode: u.DeptCode,
		Status:   u.Status,
	}
}

type RoleVO struct {
	RoleCode  string `json:"role_code"`
	RoleName  string `json:"role_name"`
	RoleKey   string `json:"role_key"`
	DataScope string `json:"data_scope"`
	Status    int    `json:"status"`
}

func BuildRoleVO(r *entity.Role) *RoleVO {
	if r == nil {
		return nil
	}
	return &RoleVO{
		RoleCode:  r.RoleCode,
		RoleName:  r.RoleName,
		RoleKey:   r.RoleKey,
		DataScope: r.DataScope,
		Status:    r.Status,
	}
}
