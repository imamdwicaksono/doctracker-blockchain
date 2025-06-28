// utils/wallet.go
package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/big"

	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
)

// GenerateMnemonic returns a new mnemonic phrase
func GenerateMnemonic() string {
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

// PrivateKeyFromMnemonic returns ECDSA key from mnemonic
func PrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, string) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal(err)
	}

	key, _ := masterKey.NewChildKey(0)
	priv := key.Key
	x, y := elliptic.P256().ScalarBaseMult(priv)
	pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	address := PublicKeyToAddress(pub)
	d := new(big.Int).SetBytes(key.Key)
	return &ecdsa.PrivateKey{D: d, PublicKey: *pub}, pub, address
}

func PublicKeyToAddress(pub *ecdsa.PublicKey) string {
	pubBytes := append(pub.X.Bytes(), pub.Y.Bytes()...)
	hash := sha256.Sum256(pubBytes)
	ripemd := ripemd160.New()
	ripemd.Write(hash[:])
	return hex.EncodeToString(ripemd.Sum(nil))
}
