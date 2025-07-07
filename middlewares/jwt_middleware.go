package middlewares

import (
	"doc-tracker/storage/redis"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Cookies("authToken") // atau dari Authorization Header

	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
	}

	// Cek token di blacklist
	isBlacklisted, err := redis.IsTokenBlacklisted(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error")
	}
	if isBlacklisted {
		return c.Status(fiber.StatusUnauthorized).SendString("Token expired or blacklisted")
	}

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

	return c.Next()
}
