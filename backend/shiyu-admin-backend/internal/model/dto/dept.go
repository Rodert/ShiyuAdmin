package dto

// CreateDeptRequest defines department creation payload.
type CreateDeptRequest struct {
	DeptCode   string `json:"dept_code" binding:"required"`
	ParentCode string `json:"parent_code"`
	DeptName   string `json:"dept_name" binding:"required"`
	Status     int    `json:"status"`
}

// UpdateDeptRequest defines department update payload.
type UpdateDeptRequest struct {
	ParentCode *string `json:"parent_code"`
	DeptName   *string `json:"dept_name"`
	Status     *int    `json:"status"`
}

