package controllers

import (
	"doc-tracker/models"
	"doc-tracker/services"
	"log"
	"os"

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
	if body.TrackerID == "" || body.Email == "" || body.Evidence == nil || *body.Evidence == "" {
		return c.Status(400).JSON(fiber.Map{"error": "tracker_id, email, and base64 evidence are required"})
	}

	checkpointAddr := services.GetCheckpointAddressByEmail(body.TrackerID, body.Email)
	if checkpointAddr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "checkpoint address not found"})
	}

	var (
		info  services.EvidenceInfo
		err   error
		useS3 = os.Getenv("S3_STORAGE") == "true"
	)

	if useS3 {
		info, err = services.SaveEvidenceBase64ToS3(body.TrackerID, checkpointAddr, body.Evidence)
	} else {
		info, err = services.SaveEvidenceFileLocal(body.TrackerID, checkpointAddr, body.Evidence)
	}
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("[Evidence] tracker=%s checkpoint=%s saved path=%s hash=%s", body.TrackerID, checkpointAddr, info.Path, info.Hash)

	if err := services.UpdateCheckpointStatus(body.TrackerID, checkpointAddr, info.Hash, info.Path); err != nil {
		return c.Status(500).JSON(fiber.Map{"error update": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":        "checkpoint complete",
		"evidence_hash": info.Hash,
		"evidence_path": info.Path,
	})
}
