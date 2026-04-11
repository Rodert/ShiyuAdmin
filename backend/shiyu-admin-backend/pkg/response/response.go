package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse defines the standard response envelope.
type APIResponse struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"error_code,omitempty"`
}

// Success sends a success response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Data:    data,
		Message: "操作成功",
	})
}

// Error sends an error response with a custom status and message.
func Error(c *gin.Context, status int, message string) {
	ErrorWithCode(c, status, "", message)
}

// ErrorWithCode sends an error response with a stable english error code and chinese message.
func ErrorWithCode(c *gin.Context, status int, errorCode, message string) {
	if status == 0 {
		status = http.StatusBadRequest
	}
	if errorCode == "" {
		errorCode = defaultErrorCode(status)
	}
	c.JSON(status, APIResponse{
		Code:      status,
		Data:      nil,
		Message:   message,
		ErrorCode: errorCode,
	})
}

func defaultErrorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusServiceUnavailable:
		return "service_unavailable"
	default:
		return "internal_server_error"
	}
}
