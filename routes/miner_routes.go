package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func MinerRoutes(app *fiber.App) {
	// Manual mining (opsional)
	app.Post("/mine", controllers.Mine)

	// Get full blockchain (debugging)
	app.Get("/chain", controllers.GetFullChain)

	// Sync chain dari peer lain
	app.Post("/sync/chain", controllers.HandleChainSync)
}
