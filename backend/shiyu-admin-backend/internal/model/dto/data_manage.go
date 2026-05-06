package dto

// DataLoginRequest defines database credential login for data monitor.
type DataLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}
