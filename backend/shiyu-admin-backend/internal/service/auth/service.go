package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"shiyu-admin-backend/internal/model/dto"
	"shiyu-admin-backend/internal/model/vo"
	"shiyu-admin-backend/internal/repository/interfaces"
	"shiyu-admin-backend/pkg/jwtutil"
)

// Service implements AuthService with repository + JWT.
type Service struct {
	repo        interfaces.AuthRepository
	jwtSecret   string
	jwtIssuer   string
	jwtExpireIn int64
}

// New creates an AuthService.
func New(repo interfaces.AuthRepository, secret, issuer string, expireIn int64) *Service {
	return &Service{
		repo:        repo,
		jwtSecret:   secret,
		jwtIssuer:   issuer,
		jwtExpireIn: expireIn,
	}
}

// Login validates credentials and returns a token.
func (s *Service) Login(ctx context.Context, req *dto.LoginRequest) (*vo.TokenVO, error) {
	user, err := s.repo.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已停用")
	}

	token, err := jwtutil.GenerateToken(s.jwtSecret, s.jwtIssuer, user.UserCode, user.Username, user.IsSuperAdmin, s.jwtExpireIn)
	if err != nil {
		return nil, err
	}
	return &vo.TokenVO{
		Token:     token,
		TokenType: "Bearer",
		ExpireIn:  s.jwtExpireIn,
	}, nil
}

