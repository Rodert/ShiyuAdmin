package jwtutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines custom JWT claims.
type Claims struct {
	CredentialStamp string `json:"credential_stamp"`
	UserCode        string `json:"user_code"`
	Username        string `json:"username"`
	IsSuperAdmin    bool   `json:"is_super_admin"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token.
func GenerateToken(secret, issuer, userCode, username string, isSuperAdmin bool, expireSeconds int64, stamps ...string) (string, error) {
	now := time.Now()
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return "", err
	}
	claims := Claims{
		UserCode:     userCode,
		Username:     username,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ID:        hex.EncodeToString(id),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	if len(stamps) > 0 {
		claims.CredentialStamp = stamps[0]
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

func CredentialStamp(secret, passwordHash string, versions ...int64) string {
	if len(versions) > 0 {
		passwordHash += ":" + strconv.FormatInt(versions[0], 10)
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(passwordHash))
	return hex.EncodeToString(h.Sum(nil))
}
