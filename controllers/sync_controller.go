package controllers

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/p2p"

	"github.com/gofiber/fiber/v2"
)

// Dipanggil saat menerima broadcast dari peer
func SyncBlock(c *fiber.Ctx) error {
	var block models.Block
	if err := c.BodyParser(&block); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid block"})
	}
	blockchain.AddBlockToChain(block)
	return c.JSON(fiber.Map{"status": "block added"})
}

// Dipanggil manual (misal tombol di UI) untuk fetch dari peer
func ManualSync(c *fiber.Ctx) error {
	peer := c.Query("peer")
	if peer == "" {
		return c.Status(400).SendString("peer is required")
	}

	block := p2p.FetchLatestBlockFrom(peer)
	blockchain.TryAddBlock(block)

	entries := p2p.FetchMempoolFrom(peer)
	for _, tx := range entries {
		mempool.AddIfNotExists(tx)
	}

	return c.SendString("Sync completed from " + peer)
}
