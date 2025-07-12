package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadECDSAPublicKey(path string) (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return ecdsaPub, nil
}

func LoadECDSAPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return priv, nil
}

func CreatePemIfNotExists(path string) error {
	// save
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("file already exists: %s", path)
	}

	privBytes, _ := x509.MarshalECPrivateKey(privKey)
	privPem := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}
	_ = os.WriteFile("data/private.pem", pem.EncodeToMemory(privPem), 0600)

	pubBytes, _ := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	_ = os.WriteFile("data/public.pem", pem.EncodeToMemory(pubPem), 0644)

	println("✅ Kunci berhasil dibuat ke data/public.pem dan private.pem")
	return nil
}
