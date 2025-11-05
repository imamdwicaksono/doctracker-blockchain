package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetMapClaims(email string) jwt.MapClaims {
	return jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
}

func ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func GenerateJWTWithClaims(claims map[string]any) (string, int64, error) {
	secret := os.Getenv("JWT_SECRET")
	expiration := time.Now().Add(time.Hour * 24)
	expUnix := expiration.Unix()

	claims["exp"] = expUnix

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, expUnix, err
}
