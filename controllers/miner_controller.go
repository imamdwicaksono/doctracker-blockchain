package controllers

import (
	"doc-tracker/blockchain"
	"doc-tracker/models"

	"github.com/gofiber/fiber/v2"
)

func GetFullChain(c *fiber.Ctx) error {
	return c.JSON(blockchain.Blockchain)
}

func HandleChainSync(c *fiber.Ctx) error {
	var incoming []blockchain.Block
	if err := c.BodyParser(&incoming); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid chain"})
	}

	if blockchain.ReplaceChain(incoming) {
		return c.JSON(fiber.Map{"message": "Blockchain replaced"})
	}
	return c.JSON(fiber.Map{"message": "No replacement needed"})
}

// (opsional) manual mining via API
func Mine(c *fiber.Ctx) error {
	// Provide appropriate values for the arguments as required by NewBlock's signature
	// Example placeholder values: index, previousHash, trackers, timestamp, nonce
	block := blockchain.NewBlock(0, "", []models.Tracker{}, 0, 0)
	blockchain.AddBlockToChain(block)
	return c.JSON(block)
}
