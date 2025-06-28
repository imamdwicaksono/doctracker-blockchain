package controllers

import (
	"doc-tracker/blockchain"

	"github.com/gofiber/fiber/v2"
)

// POST /api/sync/block
func ReceiveBlock(c *fiber.Ctx) error {
	var block blockchain.Block

	if err := c.BodyParser(&block); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid block data",
		})
	}

	lastBlock := blockchain.GetLastBlock()

	if !blockchain.IsBlockValid(block, lastBlock) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid block structure or hash",
		})
	}

	blockchain.AddBlock(block)

	return c.JSON(fiber.Map{
		"status": "Block accepted and added to chain",
	})
}

func GetChain(c *fiber.Ctx) error {
	return c.JSON(blockchain.GetAllBlocks())
}
