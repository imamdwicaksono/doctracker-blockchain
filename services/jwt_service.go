package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email     string `json:"email"`
	Address   string `json:"address,omitempty"`
	LoginWith string `json:"login_with"`           // "email" atau "tracker"
	TrackerID string `json:"tracker_id,omitempty"` // jika login via tracker
	jwt.RegisteredClaims
}

// ✅ GenerateJWT untuk email atau tracker login
func GenerateJWT(email string, loginWith string, trackerID string) (string, int64, error) {
	secret := os.Getenv("JWT_SECRET")
	expiration := time.Now().Add(time.Hour * 24) // default 24 jam
	expUnix := expiration.Unix()

	claims := JwtClaims{
		Email:     email,
		LoginWith: loginWith,
		TrackerID: trackerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, expUnix, nil
}

// ✅ Verifikasi token JWT
func VerifyJwtToken(tokenStr string) (*JwtClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
