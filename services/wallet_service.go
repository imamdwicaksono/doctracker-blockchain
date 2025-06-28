package services

import (
	"crypto/ecdsa"
	"doc-tracker/utils"
	"fmt"
	"sync"
)

type WalletInfo struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
	PublicKey  *ecdsa.PublicKey
}

var walletMap = make(map[string]WalletInfo)
var mu sync.Mutex

func GetAddressFromEmail(email string) string {
	wallet, exists := GetWalletByEmail(email)
	if !exists {
		return ""
	}
	return wallet.Address
}

func GetWalletByEmail(email string) (WalletInfo, bool) {
	mu.Lock()
	defer mu.Unlock()

	wallet, exists := walletMap[email]
	return wallet, exists
}

func GetOrCreateWallet(email string) WalletInfo {
	mu.Lock()
	defer mu.Unlock()

	if w, exists := walletMap[email]; exists {
		return w
	}

	// Gunakan mnemonic unik per email
	mnemonic := utils.GenerateMnemonic()
	privKey, pubKey, addr := utils.PrivateKeyFromMnemonic(mnemonic)

	w := WalletInfo{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Address:    addr,
	}
	walletMap[email] = w

	// Opsional: simpan mnemonic ke file atau DB jika diperlukan
	fmt.Printf("Mnemonic for %s: %s\n", email, mnemonic)
	saveMnemonicToFile(email, mnemonic)

	return w
}

func saveMnemonicToFile(email, mnemonic string) {
	// Implementasi penyimpanan mnemonic ke file atau database
	// Misalnya, simpan ke file dengan nama email.txt
	mkdir := utils.CreateDirIfNotExists("wallet/mnemonic")
	if mkdir != nil {
		fmt.Printf("Error creating directory: %v\n", mkdir)
		return
	}
	filename := fmt.Sprintf("wallet/mnemonic/%s.txt", email)
	err := utils.WriteToFile(filename, mnemonic)
	if err != nil {
		fmt.Printf("Error saving mnemonic for %s: %v\n", email, err)
	}
}
