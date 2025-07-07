package services

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email   string `json:"email"`
	Address string `json:"address"`
	jwt.RegisteredClaims
}

func VerifyJwtToken(tokenStr string) (*JwtClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
