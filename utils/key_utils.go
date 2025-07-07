package utils

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
	"os"
)

var (
	pubKeyPath  = "data/public.key"
	privKeyPath = "data/private.key"
)

// loadKeys memuat kunci dari file
func LoadKeys() (*ECDHKeyPair, error) {
	// Load private key
	privKeyBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	// Load public key
	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	curve := elliptic.P256()
	pubKey, err := DeserializePublicKey(pubKeyBytes, curve)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return &ECDHKeyPair{
		PrivateKey: new(big.Int).SetBytes(privKeyBytes),
		PublicKey:  pubKey,
	}, nil
}
