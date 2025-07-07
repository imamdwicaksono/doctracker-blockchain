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

func GenerateJWT(email string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"email": email,
		"exp":   expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return signedToken, expirationTime.Unix(), err
}

func ParseWithClaims(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
