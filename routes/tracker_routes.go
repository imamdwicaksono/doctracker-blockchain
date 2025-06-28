package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func TrackerRoutes(router fiber.Router) {
	api := router.Group("/tracker")
	api.Get("/", controllers.GetTrackers)
	api.Get("/:id", controllers.GetTrackerByID)
	api.Get("/address/:address", controllers.GetTrackersByAddress)
	api.Post("/create", controllers.CreateTracker)
}
