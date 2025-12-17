package dto

// OperationLogPageRequest defines filters and pagination for querying operation logs.
type OperationLogPageRequest struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	UserCode  string `form:"user_code"`
	Username  string `form:"username"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	Status    *int   `form:"status"`
	StartTime *int64 `form:"start_time"` // Unix timestamp seconds
	EndTime   *int64 `form:"end_time"`   // Unix timestamp seconds
}
