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
	tokenStr := c.Cookies("authToken")
	if tokenStr == "" {
		// Bisa juga ambil dari header Authorization jika cookie tidak ada
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing token",
		})
	}

	// 🔒 Cek apakah token diblacklist
	isBlacklisted, err := redis.IsTokenBlacklisted(tokenStr)
	if err != nil {
		log.Println("Redis check error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal error",
		})
	}
	if isBlacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expired or blacklisted",
		})
	}

	// 🔍 Parse dan validasi JWT
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		log.Println("Invalid JWT:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// 🧩 Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// 🧠 Simpan informasi user ke context untuk handler berikutnya
	if mode, ok := claims["mode"].(string); ok {
		c.Locals("mode", mode)
	}
	if email, ok := claims["email"].(string); ok {
		c.Locals("email", email)
	}
	if trackerID, ok := claims["tracker_id"].(string); ok {
		c.Locals("tracker_id", trackerID)
	}

	return c.Next()
}

func OnlyUserAccess(c *fiber.Ctx) error {
	mode := c.Locals("mode")
	if mode != "user" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access restricted to user accounts only",
		})
	}
	return c.Next()
}

func OnlyTrackerAccess(c *fiber.Ctx) error {
	mode := c.Locals("mode")
	if mode != "tracker" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access restricted to tracker accounts only",
		})
	}

	trackerID := c.Locals("tracker_id")
	paramID := c.Params("id")
	if trackerID != paramID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Tracker ID mismatch — not your tracker",
		})
	}

	return c.Next()
}
