package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse defines the standard response envelope.
type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
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
	if status == 0 {
		status = http.StatusBadRequest
	}
	c.JSON(status, APIResponse{
		Code:    status,
		Data:    nil,
		Message: message,
	})
}

