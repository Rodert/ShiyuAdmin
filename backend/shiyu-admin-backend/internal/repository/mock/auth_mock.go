package mock

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"shiyu-admin-backend/internal/model/entity"
	"shiyu-admin-backend/internal/repository/interfaces"
)

// AuthMockRepository is an in-memory auth repo for initial bootstrap.
type AuthMockRepository struct {
	users map[string]*entity.User
}

// NewAuthMockRepository creates a repo seeded with a default admin user.
func NewAuthMockRepository() interfaces.AuthRepository {
	// default password: admin123
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	return &AuthMockRepository{
		users: map[string]*entity.User{
			"admin": {
				ID:       1,
				UserCode: "U0001",
				Username: "admin",
				Nickname: "Administrator",
				Password: string(hash),
				Status:   1,
			},
		},
	}
}

// FindUserByUsername returns user by username.
func (m *AuthMockRepository) FindUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, nil
}

