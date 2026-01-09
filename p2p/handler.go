package p2p

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/storage"

	"github.com/gofiber/fiber/v2"
)

func FetchTrackersByAddress(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "address required"})
	}

	trackers, err := storage.Services.GetTrackersByAddress(address)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(trackers)
}

func GetLatestBlock(c *fiber.Ctx) error {
	block, err := blockchain.GetLatestBlock(storage.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(block)
}

func GetMempool(c *fiber.Ctx) error {
	return c.JSON(mempool.GetAll())
}

func GetBlocks(c *fiber.Ctx) error {
	var blocks []models.Block
	if err := storage.DB.
		Order("height asc").
		Find(&blocks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(blocks)
}

func GetBlockByHeight(c *fiber.Ctx) error {
	var block models.Block
	if err := storage.DB.
		Where("height = ?", c.Params("height")).
		First(&block).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(block)
}

func GetAudit(c *fiber.Ctx) error {
	var audit models.BlockchainAudit
	if err := storage.DB.
		Where("block_hash = ?", c.Params("blockHash")).
		First(&audit).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(audit)
}

func FetchTrackersGRPC(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Address is required"})
	}

	trackers := storage.GetTrackersByAddress(address)

	if len(trackers) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "No trackers found", "data": []models.Tracker{}})
	}

	return c.JSON(trackers)
}
