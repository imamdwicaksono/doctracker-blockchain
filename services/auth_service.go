package services

import (
	"doc-tracker/utils"
	"fmt"
)

func LoginWithMnemonic(mnemonic string) (string, error) {
	if !utils.IsValidMnemonic(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}

	_, _, address := utils.PrivateKeyFromMnemonic(mnemonic)

	return address, nil
}
