package p2p

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"

	"github.com/gofiber/fiber/v2"
)

func ReceiveBlock(c *fiber.Ctx) error {
	var block blockchain.Block
	if err := c.BodyParser(&block); err != nil {
		return c.Status(400).SendString("Invalid block format")
	}

	latest := blockchain.GetLastBlock()
	if !blockchain.IsBlockValid(block, latest) {
		return c.Status(400).SendString("Block is not valid")
	}

	// Tambahkan ke blockchain lokal
	blockchain.AddBlock(block)

	// Hapus dari mempool
	for _, tx := range block.Transactions {
		mempool.RemoveFromMempool(tx.ID)
	}

	return c.SendString("Block accepted")
}

func GetLatestBlock(c *fiber.Ctx) error {
	block := blockchain.GetLastBlock()
	return c.JSON(block)
}

func GetMempool(c *fiber.Ctx) error {
	return c.JSON(mempool.GetAll())
}
