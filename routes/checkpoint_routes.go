// routes/checkpoint_routes.go
package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterCheckpointRoutes(router fiber.Router) {
	router.Post("/checkpoint/complete", controllers.CompleteCheckpoint)
}
