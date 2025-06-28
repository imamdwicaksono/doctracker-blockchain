package controllers

import (
	"crypto/rsa"
	"doc-tracker/utils"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

type DecryptRequest struct {
	Mnemonic string `json:"mnemonic"`
	Cipher   string `json:"cipher"`
}

func DecryptNote(c *fiber.Ctx) error {
	var payload struct {
		Address string `json:"address"`
		NoteHex string `json:"note"`
		PEMKey  string `json:"pem_key"` // private key user
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).SendString("Invalid")
	}

	noteBytes, _ := hex.DecodeString(payload.NoteHex)
	privKey := utils.PrivateKeyFromPEM(payload.PEMKey)
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return c.Status(500).SendString("Invalid private key type")
	}
	decrypted, err := utils.DecryptMessage(rsaPrivKey, noteBytes)
	if err != nil {
		return c.Status(500).SendString("Decryption failed")
	}

	return c.JSON(fiber.Map{"message": string(decrypted)})
}
