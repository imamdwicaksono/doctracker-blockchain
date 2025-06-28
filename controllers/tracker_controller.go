package controllers

import (
	"doc-tracker/models"
	"doc-tracker/services"

	"github.com/gofiber/fiber/v2"
)

func CreateTracker(c *fiber.Ctx) error {
	var input models.Tracker
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.Type == "" || input.Privacy == "" || input.Creator == "" || len(input.Checkpoints) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	tracker, err := services.CreateTracker(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create tracker"})
	}

	return c.JSON(tracker)
}

func GetTrackers(c *fiber.Ctx) error {
	trackers, err := services.GetTrackers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch trackers"})
	}
	if len(trackers) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "No trackers found", "data": []models.Tracker{}})
	}
	// Filter out trackers with status "progress"
	// var filteredTrackers []models.Tracker
	// for _, t := range trackers {
	// 	if t.Status != "progress" {
	// 		filteredTrackers = append(filteredTrackers, t)
	// 	}
	// }
	return c.JSON(trackers)
}

func GetTrackerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tracker ID is required"})
	}

	tracker, err := services.GetTrackerByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tracker not found"})
	}

	return c.JSON(tracker)
}

func GetTrackersByAddress(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Address is required"})
	}

	trackers, err := services.GetTrackersByAddress(address)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch trackers by address"})
	}

	if len(trackers) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "No trackers found", "data": []models.Tracker{}})
	}

	return c.JSON(trackers)
}
