package apperrors

import (
	"errors"
	"net/http"
	"strings"
)

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, code, message string) error {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func NewConflict(code, message string) error {
	return New(http.StatusConflict, code, message)
}

func IsConflict(err error) bool {
	var target *AppError
	return errors.As(err, &target) && target.Status == http.StatusConflict
}

func Extract(err error) (*AppError, bool) {
	var target *AppError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

func WrapUniqueConstraint(err error, code, message string) error {
	if IsUniqueConstraint(err) {
		return NewConflict(code, message)
	}
	return err
}

func IsUniqueConstraint(err error) bool {
	if err == nil {
		return false
	}
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "duplicate entry") ||
		strings.Contains(lower, "duplicate key value") ||
		strings.Contains(lower, "duplicated key") ||
		strings.Contains(lower, "unique constraint") ||
		strings.Contains(lower, "unique failed") ||
		strings.Contains(lower, "error 1062")
}
