package routes

import (
	"doc-tracker/controllers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
	auth.Post("/request-otp", controllers.SendOtp)
	auth.Post("/verify-otp", controllers.VerifyOtp)
}

func SetupAuthProtectedRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Get("/qr/:address", controllers.GetQR)

	meauth := router.Group("/auth")
	meauth.Use(limiter.New(limiter.Config{Max: 10000, Expiration: time.Minute}))
	meauth.Get("/me", controllers.AuthMe)
}
