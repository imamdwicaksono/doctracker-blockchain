package middlewares

import (
	"doc-tracker/storage/redis"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Ambil dari Authorization Header
	authHeader := c.Get("Authorization")
	var tokenStr string

	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Fallback: ambil dari cookie
		tokenStr = c.Cookies("authToken")
	}

	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
	}

	// Cek blacklist di Redis
	isBlacklisted, err := redis.IsTokenBlacklisted(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}
	if isBlacklisted {
		return c.Status(fiber.StatusUnauthorized).SendString("Token expired or blacklisted")
	}

	// Parse dan verifikasi token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		log.Println("Invalid JWT:", err)
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
	}

	// Optional: extract claim (email, address, dll)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if email, ok := claims["email"].(string); ok {
			c.Locals("user_email", email)
		}
	}

	return c.Next()
}
