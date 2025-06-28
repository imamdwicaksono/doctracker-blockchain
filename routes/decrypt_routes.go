package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterDecryptRoutes(router fiber.Router) {
	router.Post("/decrypt-note", controllers.DecryptNote)
}
