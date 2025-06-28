package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterEvidenceRoutes(router fiber.Router) {
	router.Post("/upload", controllers.UploadEvidence)
}
