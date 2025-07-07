package routes

import (
	"doc-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func TrackerRoutes(router fiber.Router) {
	api := router.Group("/trackers")
	api.Get("/", controllers.GetTrackers)

	apiTracker := router.Group("/tracker")
	apiTracker.Get("/:id", controllers.GetTrackerByID)
	apiTracker.Get("/address/:address", controllers.GetTrackersByAddress)
	apiTracker.Post("/create", controllers.CreateTracker)
	apiTracker.Get("/summary/:email", controllers.GetTrackerSummary)
}
