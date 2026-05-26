package interfaces

import "context"

// UserDataScope describes the rows a user may see in user-owned datasets.
type UserDataScope struct {
	All       bool
	UserCode  string
	DeptCodes []string
}

// DataScopeService resolves role data-scope settings into reusable filters.
type DataScopeService interface {
	ResolveUserScope(ctx context.Context, userCode string, isSuperAdmin bool) (*UserDataScope, error)
}
