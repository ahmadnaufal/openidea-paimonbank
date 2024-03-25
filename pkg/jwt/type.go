package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUser struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func BuildJWTClaims(user JWTUser, expireDuration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"userId": user.UserID,
		"name":   user.Name,
		"email":  user.Email,
		"exp":    jwt.NewNumericDate(time.Now().Add(expireDuration)),
	}
}
