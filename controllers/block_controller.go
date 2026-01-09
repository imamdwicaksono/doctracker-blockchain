package controllers

import (
	"doc-tracker/models"
	"doc-tracker/storage"

	"github.com/gofiber/fiber/v2"
)

// GET /api/blocks
// List semua block audit (DB-first)
func GetBlocks(c *fiber.Ctx) error {
	var blocks []models.Block

	if err := storage.DB.
		Order("height asc").
		Find(&blocks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(blocks)
}

// GET /api/blocks/:height
// Ambil block berdasarkan height
func GetBlockByHeight(c *fiber.Ctx) error {
	height := c.Params("height")

	var block models.Block
	if err := storage.DB.
		Where("height = ?", height).
		First(&block).Error; err != nil {

		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "block not found"})
	}

	return c.JSON(block)
}

func SyncBlock(c *fiber.Ctx) error {
	var block models.Block

	if err := c.BodyParser(&block); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid block format"})
	}

	// Simpan ke DB
	if err := storage.DB.Create(&block).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "block synced"})
}
