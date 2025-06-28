package controllers

import (
	"doc-tracker/services"
	"doc-tracker/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Mnemonic string `json:"mnemonic"`
}

func Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	address, err := services.LoginWithMnemonic(input.Mnemonic)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"address": address,
	})
}

func GetQR(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing address")
	}

	png, err := utils.GenerateQRCode(address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate QR")
	}

	c.Type("png")
	return c.Send(png)
}
