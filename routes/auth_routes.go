package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {
	router.Post("/login", controllers.Login)
	router.Get("/qr/:address", controllers.GetQR)
}
