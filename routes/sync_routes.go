package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func SyncRoutes(app *fiber.App) {
	app.Post("/p2p/block", controllers.SyncBlock)   // receive push
	app.Get("/sync/manual", controllers.ManualSync) // manual fetch
}
