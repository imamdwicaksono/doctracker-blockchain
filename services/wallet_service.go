package services

import (
	"crypto/ecdsa"
	"doc-tracker/utils"
	"fmt"
	"os"
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

	wallet := GetWalletFromPem(email)
	if wallet.PrivateKey == nil || wallet.PublicKey == nil || wallet.Address == "" {
		fmt.Printf("Wallet for %s is not properly initialized\n", email)
		return ""
	}
	return wallet.Address
}

func GetWalletFromPem(email string) WalletInfo {
	mu.Lock()
	defer mu.Unlock()

	wallet, err := os.ReadFile("wallet/mnemonic/" + email + ".txt")
	if err != nil {
		fmt.Printf("Error reading wallet file for %s: %v\n", email, err)
		return WalletInfo{}
	}

	// Jika wallet tidak ada, buat baru
	privateKey, publicKey, address := utils.PrivateKeyFromMnemonic(string(wallet))

	return WalletInfo{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}

func GetWalletByEmail(email string) (WalletInfo, bool) {
	mu.Lock()

	defer mu.Unlock()

	wallet := GetOrCreateWallet(email) // Pastikan wallet sudah ada atau dibuat
	if wallet.PrivateKey == nil || wallet.PublicKey == nil || wallet.Address == "" {
		fmt.Printf("Wallet for %s is not properly initialized\n", email)
		return WalletInfo{}, false
	}
	return wallet, true
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
	fmt.Printf("Mnemonic for %s saved to %s\n", email, filename)
}
