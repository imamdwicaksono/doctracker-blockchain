package controllers

import (
	"crypto/sha256"
	"doc-tracker/services"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func UploadEvidence(c *fiber.Ctx) error {
	trackerID := c.FormValue("tracker_id")
	checkpointAddr := c.FormValue("checkpoint_address")

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File not found")
	}

	// Open file stream
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Can't open file")
	}
	defer src.Close()

	// Hash the content
	hasher := sha256.New()
	io.Copy(hasher, src)

	hash := hex.EncodeToString(hasher.Sum(nil))

	// Reset src reader
	src.Seek(0, io.SeekStart)

	// Save file
	dir := fmt.Sprintf("./storage/evidence/%s", trackerID)
	os.MkdirAll(dir, 0755)
	dstPath := filepath.Join(dir, checkpointAddr+filepath.Ext(file.Filename))

	dst, err := os.Create(dstPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file")
	}
	defer dst.Close()
	io.Copy(dst, src)

	// Update checkpoint status
	err = services.UpdateCheckpointStatus(trackerID, checkpointAddr, hash, dstPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message":       "Evidence uploaded",
		"evidence_hash": hash,
		"evidence_path": dstPath,
	})
}
