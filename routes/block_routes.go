package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func BlockRoutes(app fiber.Router) {
	app.Post("/sync/block", controllers.SyncBlock)
	app.Get("/blocks", controllers.GetBlocks)
}
