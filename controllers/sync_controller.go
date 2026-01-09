package controllers

import (
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/storage"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

// POST /api/sync/tracker
func SyncTracker(c *fiber.Ctx) error {
	var tracker models.Tracker

	if err := c.BodyParser(&tracker); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid tracker"})
	}

	// Simpan ke DB (idempotent)
	if err := storage.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).
		Create(&tracker).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Masukkan ke mempool cache
	mempool.AddIfNotExists(tracker)

	return c.JSON(fiber.Map{"status": "tracker synced"})
}

// POST /api/sync/mempool
func SyncMempool(c *fiber.Ctx) error {
	var trackers []models.Tracker

	if err := c.BodyParser(&trackers); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid mempool payload"})
	}

	for _, t := range trackers {
		// DB first
		_ = storage.DB.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).
			Create(&t).Error

		mempool.AddIfNotExists(t)
	}

	return c.JSON(fiber.Map{"status": "mempool synced"})
}
