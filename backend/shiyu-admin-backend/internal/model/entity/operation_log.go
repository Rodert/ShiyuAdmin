package entity

import "time"

// OperationLog represents a user operation audit log.
type OperationLog struct {
	// ID 主键ID，自增
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID，自增"`
	// UserCode 用户编码，关联用户表
	UserCode string `json:"user_code" gorm:"size:32;index;comment:用户编码，关联用户表"`
	// Username 用户名，冗余存储便于查询
	Username string `json:"username" gorm:"size:64;comment:用户名，冗余存储便于查询"`
	// Module 业务模块，如 system-user、system-role
	Module string `json:"module" gorm:"size:64;comment:业务模块，如 system-user、system-role"`
	// Action 操作动作，如 create、update、delete
	Action string `json:"action" gorm:"size:64;comment:操作动作，如 create、update、delete"`
	// Method HTTP 方法
	Method string `json:"method" gorm:"size:16;comment:HTTP 方法"`
	// Path 请求路径
	Path string `json:"path" gorm:"size:255;comment:请求路径"`
	// IP 客户端 IP
	IP string `json:"ip" gorm:"size:64;comment:客户端 IP"`
	// Status 操作状态：1=成功，0=失败
	Status int `json:"status" gorm:"comment:操作状态：1=成功，0=失败"`
	// ErrorMsg 错误信息摘要（截断）
	ErrorMsg string `json:"error_msg" gorm:"size:512;comment:错误信息摘要（截断）"`
	// LatencyMs 耗时（毫秒）
	LatencyMs int64 `json:"latency_ms" gorm:"comment:耗时（毫秒）"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
}

// TableName 自定义表名
func (OperationLog) TableName() string { return "sys_operation_logs" }
