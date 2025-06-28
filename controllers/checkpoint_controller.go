// controllers/checkpoint_controller.go
package controllers

import (
	"doc-tracker/services"

	"github.com/gofiber/fiber/v2"
)

func CompleteCheckpoint(c *fiber.Ctx) error {
	var body struct {
		TrackerID string  `json:"tracker_id"`
		Email     string  `json:"email"`
		Note      *string `json:"note"`
		// Evidence is a base64 encoded string or file path
		Evidence *string `json:"evidence"`
	}

	// FIX: hanya parsing body, tidak return langsung
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	return c.JSON(body)

	if body.TrackerID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "tracker_id is required"})
	}

	if body.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "email is required"})
	}

	checkpointAddr := services.GetAddressFromEmail(body.Email)
	if checkpointAddr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid checkpoint email"})
	}

	if body.Evidence == nil || *body.Evidence == "" {
		return c.Status(400).JSON(fiber.Map{"error": "evidence file or base64 string is required"})
	}

	formFile, err := services.SaveBase64Evidence(body.TrackerID, checkpointAddr, body.Evidence)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid base64 file"})
	}

	info, err := services.SaveEvidenceFile(body.TrackerID, checkpointAddr, formFile)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	err = services.UpdateCheckpointStatus(body.TrackerID, checkpointAddr, info.Hash, info.Path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":        "checkpoint complete",
		"evidence_hash": info.Hash,
		"evidence_path": info.Path,
	})
}
