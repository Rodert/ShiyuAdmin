package dto

// SetRoleDeptsRequest defines role-dept assignment payload.
type SetRoleDeptsRequest struct {
	DeptCodes []string `json:"dept_codes" binding:"required"`
}

