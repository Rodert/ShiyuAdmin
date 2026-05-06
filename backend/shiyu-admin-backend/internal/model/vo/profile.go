package vo

// ProfileVO is the current user's profile response.
type ProfileVO struct {
	UserCode     string          `json:"user_code"`
	Username     string          `json:"username"`
	Nickname     string          `json:"nickname"`
	Email        string          `json:"email"`
	Phone        string          `json:"phone"`
	Avatar       string          `json:"avatar"`
	DeptCode     string          `json:"dept_code"`
	Status       int             `json:"status"`
	IsSuperAdmin bool            `json:"is_super_admin"`
	Dept         *ProfileDeptVO  `json:"dept,omitempty"`
	Roles        []ProfileRoleVO `json:"roles"`
}

type ProfileDeptVO struct {
	DeptCode string `json:"dept_code"`
	DeptName string `json:"dept_name"`
}

type ProfileRoleVO struct {
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
	RoleKey  string `json:"role_key"`
}
