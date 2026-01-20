package controllers

import (
	"doc-tracker/models"
	"doc-tracker/services"
	"doc-tracker/storage/jwt"
	"doc-tracker/storage/redis"
	"doc-tracker/utils"
	"fmt"
	"math/rand"
	"os"
	"strings"
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
		fmt.Println("❌ BodyParser failed:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	identifier := req.Email // digunakan untuk dua fungsi: email atau tracker ID
	if identifier == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Identifier (email or tracker ID) is required")
	}

	loginWith := "email"
	trackerID := ""

	// === Jika input adalah Tracker ID ===
	if strings.HasPrefix(strings.ToUpper(identifier), "TRK-") {
		trackerID = identifier
		fmt.Println("🔎 Detected Tracker ID:", trackerID)

		// Dapatkan data tracker dari database
		tracker, err := services.GetTrackerByID(trackerID)
		if err != nil {
			fmt.Println("❌ Error get tracker:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  500,
				"message": "Failed to fetch tracker",
			})
		}
		if tracker.ID == "" {
			fmt.Println("❌ Tracker not found:", trackerID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  404,
				"message": "No tracker found with that ID",
			})
		}

		// Gunakan email creator tracker untuk kirim OTP
		if tracker.Creator == "" {
			fmt.Println("❌ Tracker has no creator email")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  400,
				"message": "Tracker creator email not found",
			})
		}

		identifier = tracker.Creator
		loginWith = "tracker"
		fmt.Printf("📩 Tracker login: sending OTP to creator email %s (tracker %s)\n", identifier, trackerID)
	}

	// === Validasi email format sederhana ===
	if !strings.Contains(identifier, "@") {
		fmt.Println("❌ Invalid email format:", identifier)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}

	// === Generate dan kirim OTP ===
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	fmt.Printf("📨 Sending OTP %s to %s\n", otp, identifier)

	// Simpan OTP ke Redis / memory
	if err := redis.StoreOtpInMemoryOrRedis(identifier, otp); err != nil {
		fmt.Println("❌ Failed storing OTP:", err)
	}

	// Kirim email OTP
	if err := utils.SendEmailOTP(identifier, otp, "resources/emails/request-otp.html"); err != nil {
		fmt.Println("❌ Failed to send email:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send OTP email")
	}

	fmt.Println("✅ OTP sent successfully")

	// === Response ke frontend ===
	resp := fiber.Map{
		"status":     200,
		"login_with": loginWith,
		"message":    fmt.Sprintf("OTP sent to %s successfully", identifier),
	}

	// Jika login dengan tracker, tambahkan tracker_id agar bisa redirect
	if loginWith == "tracker" {
		resp["tracker_id"] = trackerID
	}

	return c.JSON(resp)
}

func VerifyOtp(c *fiber.Ctx) error {
	var req models.VerifyOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	identifier := strings.TrimSpace(req.Email) // Bisa email atau TRK-xxxx
	isTrackerLogin := strings.HasPrefix(strings.ToUpper(identifier), "TRK-")

	var emailForOtp string
	var trackerID string

	// --- [1] Deteksi mode login ---
	if isTrackerLogin {
		trackerID = identifier

		// Ambil data tracker untuk dapatkan email pemilik
		tracker, err := services.GetTrackerByID(trackerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tracker"})
		}
		if tracker.ID == "" {
			return c.Status(404).JSON(fiber.Map{"message": "No tracker found", "data": []models.Tracker{}})
		}

		emailForOtp = tracker.Creator
		if emailForOtp == "" {
			return fiber.NewError(fiber.StatusNotFound, "Tracker not found or email not linked")
		}
	} else {
		emailForOtp = identifier
	}

	// --- [2] Validasi OTP ---
	expectedOtp := redis.GetOtpFromMemoryOrRedis(emailForOtp)
	if expectedOtp == "" || req.Otp != expectedOtp {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid OTP")
	}

	// Hapus OTP setelah dipakai
	redis.Client.Del(redis.Ctx, "otp:"+emailForOtp)
	fmt.Println("✅ OTP verified successfully, removing from cache")

	// --- [3] Generate JWT ---
	var (
		token   string
		expUnix int64
		err     error
	)

	if isTrackerLogin {
		token, expUnix, err = services.GenerateJWT(emailForOtp, "tracker", trackerID)
	} else {
		token, expUnix, err = services.GenerateJWT(emailForOtp, "email", "")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	fmt.Println("✅ JWT token generated successfully")

	// --- [4] Set cookie auth ---
	maxAge := 86400 // default 1 hari
	if v := os.Getenv("COOKIE_MAX_AGE"); v != "" {
		fmt.Sscanf(v, "%d", &maxAge)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    token,
		HTTPOnly: true,
		Secure:   os.Getenv("COOKIE_SECURE") == "true",
		Path:     os.Getenv("COOKIE_PATH"),
		MaxAge:   maxAge,
		SameSite: os.Getenv("COOKIE_SAMESITE"),
		Domain:   os.Getenv("COOKIE_DOMAIN_NAME"),
	})

	fmt.Printf("✅ Cookie set (mode=%s)\n", func() string {
		if isTrackerLogin {
			return "tracker"
		}
		return "user"
	}())

	// --- [5] Response sukses ---
	return c.JSON(fiber.Map{
		"status":     200,
		"message":    "OTP verified successfully",
		"token":      token,
		"email":      emailForOtp,
		"tracker_id": trackerID,
		"mode": func() string {
			if isTrackerLogin {
				return "tracker"
			}
			return "user"
		}(),
		"exp": expUnix,
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

	// Clear cookie
	c.ClearCookie("authToken") // Nama cookie authToken
	// Hapus cookie (opsional jika pakai header)
	maxAge := 0
	if v := os.Getenv("COOKIE_MAX_AGE"); v != "" {
		fmt.Sscanf(v, "%d", &maxAge)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "authToken",
		Value:    "",
		HTTPOnly: true,
		Secure:   os.Getenv("COOKIE_SECURE") == "true", // ⬅️ WAJIB true jika pakai SameSite=None
		Path:     os.Getenv("COOKIE_PATH"),
		MaxAge:   maxAge,                          // ⬅️ WAJIB sesuai dengan TTL token
		SameSite: os.Getenv("COOKIE_SAMESITE"),    // ⬅️ WAJIB "None" agar bisa cross-domain
		Domain:   os.Getenv("COOKIE_DOMAIN_NAME"), // ⬅️ optional tapi bisa bantu konsisten
	})

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Logged out successfully",
	})
}

func AuthMe(c *fiber.Ctx) error {
	tokenStr := c.Cookies("authToken")

	// Fallback: cari di Authorization header
	if tokenStr == "" {
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = authHeader[7:]
		}
	}

	if tokenStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}

	// Verifikasi token JWT
	claims, err := services.VerifyJwtToken(tokenStr)
	if err != nil {
		fmt.Println("❌ Invalid token:", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Deteksi mode login berdasarkan claims
	loginWith := "email"
	trackerID := ""

	if claims.LoginWith == "tracker" && claims.TrackerID != "" {
		loginWith = "tracker"
		trackerID = claims.TrackerID
	}

	resp := fiber.Map{
		"status":     200,
		"login_with": loginWith,
		"email":      claims.Email,
		"address":    claims.Address,
	}

	if trackerID != "" {
		resp["tracker_id"] = trackerID
	}

	return c.JSON(resp)
}
