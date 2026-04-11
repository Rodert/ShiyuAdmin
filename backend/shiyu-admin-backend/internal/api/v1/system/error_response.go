package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shiyu-admin-backend/internal/apperrors"
	"shiyu-admin-backend/pkg/response"
)

func writeServiceError(c *gin.Context, err error) {
	if appErr, ok := apperrors.Extract(err); ok {
		response.ErrorWithCode(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}
	response.ErrorWithCode(c, http.StatusInternalServerError, "internal_server_error", err.Error())
}
