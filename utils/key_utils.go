package utils

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
)

func LoadKeysFromEnv() (*ECDHKeyPair, error) {
	privHex := os.Getenv("ECDH_PRIVATE_KEY")
	pubHex := os.Getenv("ECDH_PUBLIC_KEY")

	if privHex == "" || pubHex == "" {
		return nil, fmt.Errorf("ECDH_PRIVATE_KEY or ECDH_PUBLIC_KEY not set")
	}

	privBytes, err := HexToBytes(privHex)
	if err != nil {
		return nil, err
	}

	pubBytes, err := HexToBytes(pubHex)
	if err != nil {
		return nil, err
	}

	curve := elliptic.P256()

	d := new(big.Int).SetBytes(privBytes)
	d.Mod(d, curve.Params().N)

	pubKey, err := DeserializePublicKey(pubBytes, curve)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	if !curve.IsOnCurve(pubKey.X, pubKey.Y) {
		return nil, fmt.Errorf("public key not on curve")
	}

	return &ECDHKeyPair{
		PrivateKey: d,
		PublicKey:  pubKey,
	}, nil
}

func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		return nil, fmt.Errorf("invalid hex string length")
	}
	return hex.DecodeString(hexStr)
}
