package vo

// TokenVO represents login response payload.
type TokenVO struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"` // e.g., Bearer
	ExpireIn  int64  `json:"expire_in"`  // seconds
}

