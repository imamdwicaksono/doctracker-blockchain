package services

import (
	"doc-tracker/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoginWithMnemonic(mnemonic string) (string, error) {
	if !utils.IsValidMnemonic(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}

	_, _, address := utils.PrivateKeyFromMnemonic(mnemonic)

	return address, nil
}

func GetLoginEmail(c *fiber.Ctx) (string, error) {
	tokenStr := c.Cookies("authToken")

	// Fallback: cari di Authorization header
	if tokenStr == "" {
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		}
	}

	if tokenStr == "" {
		return "", fmt.Errorf("missing token")
	}

	claims, err := VerifyJwtToken(tokenStr)
	if err != nil {
		return "", err
	}

	return claims.Email, nil
}
