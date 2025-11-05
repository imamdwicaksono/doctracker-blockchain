package controllers

import (
	"doc-tracker/models"
	"doc-tracker/services"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CreateTracker godoc
// @Summary     Create a new document tracker
// @Description Create a tracker with initial checkpoints
// @Tags        Trackers
// @Accept      json
// @Produce     json
// @Param       tracker body models.Tracker true "Tracker payload"
// @Success     200 {object} models.Tracker
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /trackers [post]
func CreateTracker(c *fiber.Ctx) error {
	var input models.Tracker
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.Type == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Type is required"})
	}

	if len(input.Checkpoints) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "At least one checkpoint is required"})
	}

	first := input.Checkpoints[0]

	// ✅ Validasi email / emails
	hasEmail := false
	if len(first.Emails) > 0 {
		for _, e := range first.Emails {
			if strings.TrimSpace(e) != "" {
				hasEmail = true
				break
			}
		}
	} else if strings.TrimSpace(first.Email) != "" {
		hasEmail = true
	}

	if !hasEmail {
		return c.Status(400).JSON(fiber.Map{"error": "Email is required for the first checkpoint"})
	}

	// ✅ Proses lanjut
	data, err := services.CreateTracker(input)
	if err != nil {
		fmt.Printf("CreateTracker error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create tracker"})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  201,
		"message": "Tracker created successfully",
		"data":    data,
	})
}

// GetTrackers godoc
// @Summary     Get all document trackers
// @Description Retrieve all trackers
// @Tags        Trackers
// @Produce     json
// @Success     200 {array} models.Tracker
// @Failure     500 {object} map[string]string
// @Router      /trackers [get]
func GetTrackers(c *fiber.Ctx) error {
	trackers, err := services.GetTrackers(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch trackers"})
	}
	if len(trackers) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "No trackers found", "data": []models.Tracker{}})
	}
	return c.JSON(trackers)
}

// GetTrackerByID godoc
// @Summary     Get a tracker by ID
// @Description Retrieve tracker using its ID
// @Tags        Trackers
// @Produce     json
// @Param       id path string true "Tracker ID"
// @Success     200 {object} models.Tracker
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /trackers/{id} [get]
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

// GetTrackersByAddress godoc
// @Summary     Get trackers by wallet address
// @Description Retrieve all trackers initiated by a specific address
// @Tags        Trackers
// @Produce     json
// @Param       address path string true "Wallet address"
// @Success     200 {array} models.Tracker
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]string
// @Router      /trackers/address/{address} [get]
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

func GetTrackerSummary(c *fiber.Ctx) error {
	email := c.Params("email")
	summary, err := services.GetTrackerSummary(email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tracker summary"})
	}
	return c.JSON(fiber.Map{"data": summary})
}
