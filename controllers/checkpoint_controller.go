package controllers

import (
	"doc-tracker/models"
	"doc-tracker/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// CompleteCheckpoint godoc
// @Summary     Complete a checkpoint with evidence
// @Description Mark a checkpoint as completed by uploading evidence
// @Tags        Checkpoints
// @Accept      json
// @Produce     json
// @Param checkpoint body models.CheckpointStatusInput true "Checkpoint Status Input"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /checkpoint/complete [post]
func CompleteCheckpoint(c *fiber.Ctx) error {
	var body models.CheckpointStatusInput

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if body.TrackerID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "tracker_id is required"})
	}

	if body.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "email is required"})
	}

	if *body.Evidence == "" {
		return c.Status(400).JSON(fiber.Map{"error": "base64 string is required"})
	}
	fmt.Printf("Received base64 evidence for tracker %s at checkpoint %s\n", body.TrackerID, body.Email)

	// formFile, err := services.SaveBase64Evidence(body.TrackerID, checkpointAddr, body.Evidence)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"error": "invalid base64 file"})
	// }
	// fmt.Printf("Evidence file saved for tracker %s at checkpoint %s: %s\n", body.TrackerID, checkpointAddr, formFile)
	checkPointAddr := services.GetCheckpointAddressByEmail(body.TrackerID, body.Email)
	if checkPointAddr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "checkpoint address not found"})
	}

	info, err := services.SaveEvidenceFile(body.TrackerID, checkPointAddr, body.Evidence)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error save evidence": err.Error()})
	}
	fmt.Printf("Evidence file processed for tracker %s at checkpoint %s: %s, hash: %s\n", body.TrackerID, checkPointAddr, info.Path, info.Hash)

	err = services.UpdateCheckpointStatus(body.TrackerID, checkPointAddr, info.Hash, info.Path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error update": err.Error()})
	}
	fmt.Printf("Checkpoint status updated for tracker %s at checkpoint %s with hash %s\n", body.TrackerID, checkPointAddr, info.Hash)
	fmt.Printf("Evidence file path: %s\n", info.Path)
	fmt.Printf("Evidence file hash: %s\n", info.Hash)
	return c.JSON(fiber.Map{
		"status":        "checkpoint complete",
		"evidence_hash": info.Hash,
		"evidence_path": info.Path,
	})
}
