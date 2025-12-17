package vo

import "shiyu-admin-backend/internal/model/entity"

// OperationLogVO represents operation log data for API responses.
type OperationLogVO struct {
	ID        int64  `json:"id"`
	UserCode  string `json:"user_code"`
	Username  string `json:"username"`
	Module    string `json:"module"`
	Action    string `json:"action"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	IP        string `json:"ip"`
	Status    int    `json:"status"`
	ErrorMsg  string `json:"error_msg"`
	LatencyMs int64  `json:"latency_ms"`
	CreatedAt int64  `json:"created_at"`
}

// BuildOperationLogVO converts entity to VO.
func BuildOperationLogVO(l *entity.OperationLog) *OperationLogVO {
	if l == nil {
		return nil
	}
	return &OperationLogVO{
		ID:        l.ID,
		UserCode:  l.UserCode,
		Username:  l.Username,
		Module:    l.Module,
		Action:    l.Action,
		Method:    l.Method,
		Path:      l.Path,
		IP:        l.IP,
		Status:    l.Status,
		ErrorMsg:  l.ErrorMsg,
		LatencyMs: l.LatencyMs,
		CreatedAt: l.CreatedAt.Unix(),
	}
}
