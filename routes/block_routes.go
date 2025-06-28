package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func BlockRoutes(app fiber.Router) {
	app.Post("/sync/block", controllers.ReceiveBlock)
	app.Get("/blocks", controllers.GetChain)
}
