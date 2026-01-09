package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MinerRoutes(app *fiber.App, db *gorm.DB) {
	// Manual mining (opsional)
	app.Post("/mine", func(c *fiber.Ctx) error {
		return controllers.Mine(c, db)
	})

	// Get full blockchain (debugging)
	app.Get("/chain", controllers.GetBlocks)

	// Sync chain dari peer lain
	app.Post("/sync/chain", controllers.SyncBlock)
}
