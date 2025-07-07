package controllers

import (
	"doc-tracker/models"
	"doc-tracker/services"
	"doc-tracker/storage/jwt"
	"doc-tracker/storage/redis"
	"doc-tracker/utils"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Mnemonic string `json:"mnemonic"`
}

func Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	address, err := services.LoginWithMnemonic(input.Mnemonic)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"address": address,
	})
}

func SendOtp(c *fiber.Ctx) error {
	var req models.OtpRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	redis.StoreOtpInMemoryOrRedis(req.Email, otp)

	// Kirim ke email (SMTP)
	if err := utils.SendEmailOTP(req.Email, otp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send email")
	}

	return c.JSON(fiber.Map{"status": 200, "message": "OTP sent successfully"})
}

func VerifyOtp(c *fiber.Ctx) error {
	var req models.VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	expectedOtp := redis.GetOtpFromMemoryOrRedis(req.Email)
	if expectedOtp == "" || req.Otp != expectedOtp {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid OTP")
	}

	// Hapus OTP setelah digunakan
	redis.Client.Del(redis.Ctx, "otp:"+req.Email)

	// Buat JWT token
	token, expUnix, err := jwt.GenerateJWT(req.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    token,
		HTTPOnly: true,
		Secure:   false, // ubah jadi true kalau pakai HTTPS
		Path:     "/",
		MaxAge:   86400,
		SameSite: "Lax", // atau "Strict" / "None" tergantung skenario
	})

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "OTP verified successfully",
		"token":   token,
		"email":   req.Email,
		"exp":     expUnix,
	})
}

func GetQR(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing address")
	}

	png, err := utils.GenerateQRCode(address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate QR")
	}

	c.Type("png")
	return c.Send(png)
}

func Logout(c *fiber.Ctx) error {
	token := c.Cookies("authToken")

	// Fallback: ambil dari Authorization header
	if token == "" {
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing token",
		})
	}

	// Optional: parsing untuk TTL
	claims := jwt.GetMapClaims("")
	parsed, err := jwt.ParseWithClaims(token, claims)
	if err == nil && parsed.Valid {
		if expUnix, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(expUnix), 0)
			ttl := time.Until(expTime)
			_ = redis.BlacklistToken(token, ttl)
		}
	}

	// Hapus cookie (opsional jika pakai header)
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func AuthMe(c *fiber.Ctx) error {
	tokenStr := c.Cookies("authToken")

	// Fallback: cari di Authorization header
	if tokenStr == "" {
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		}
	}

	if tokenStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}

	claims, err := services.VerifyJwtToken(tokenStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(fiber.Map{
		"email":   claims.Email,
		"address": claims.Address,
	})
}
