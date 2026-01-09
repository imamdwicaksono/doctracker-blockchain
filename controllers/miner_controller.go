package controllers

import (
	"doc-tracker/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Mine(c *fiber.Ctx, db *gorm.DB) error {
	go func() {
		for {
			err := services.MineOnce(db)
			if err != nil {
				break
			}
		}
	}()
	return c.JSON(fiber.Map{"status": "mining started"})
}
