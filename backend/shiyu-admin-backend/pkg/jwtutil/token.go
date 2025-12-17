package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines custom JWT claims.
type Claims struct {
	UserCode string `json:"user_code"`
	Username string `json:"username"`
	IsSuperAdmin bool `json:"is_super_admin"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token.
func GenerateToken(secret, issuer, userCode, username string, isSuperAdmin bool, expireSeconds int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserCode:     userCode,
		Username:     username,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates and parses JWT token string.
func ParseToken(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
