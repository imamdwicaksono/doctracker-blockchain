package utils

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"

	ecies "github.com/ecies/go/v2"
)

// ConvertAndEncryptJSONToBin reads a JSON file, encrypts it, and saves as .bin
func ConvertAndEncryptJSONToBin(jsonPath, publicKeyPath, outputBinPath string) error {
	// Load public key
	pubECDSA, err := LoadPEMPub(publicKeyPath)
	if err != nil {
		return errors.New("Load Pem: " + err.Error())
	}
	pub := ECDSAPublicKeyToECIES(pubECDSA)

	// Read JSON file
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return errors.New("Read JSON: " + err.Error())
	}

	// Encrypt content
	encData, err := ecies.Encrypt(pub, data)
	if err != nil {
		return err
	}

	// Write to .bin
	return os.WriteFile(outputBinPath, encData, 0644)
}

// Enkripsi file menggunakan public key dari file PEM
func EncryptFileWithPEM(publicKeyPath, inputPath, outputPath string) error {
	pubECDSA, err := LoadPEMPub(publicKeyPath)
	if err != nil {
		return fmt.Errorf("LoadPEMPub error: %w", err)
	}

	pub := ECDSAPublicKeyToECIES(pubECDSA)

	plainData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	encData, err := ecies.Encrypt(pub, plainData)
	if err != nil {
		return fmt.Errorf("encrypt error: %w", err)
	}

	return os.WriteFile(outputPath, encData, 0644)
}

// Dekripsi file menggunakan private key dari file PEM
func DecryptFileWithPEM(privateKeyPath, inputPath, outputPath string) error {
	privECDSA, err := LoadPEMKey(privateKeyPath)
	if err != nil {
		return fmt.Errorf("LoadPEMKey error: %w", err)
	}

	priv := ECDSAPrivateKeyToECIES(privECDSA)

	encData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read encrypted input: %w", err)
	}

	plainData, err := ecies.Decrypt(priv, encData)
	if err != nil {
		return fmt.Errorf("decrypt error: %w", err)
	}

	return os.WriteFile(outputPath, plainData, 0644)
}

func ECDSAPublicKeyToECIES(pub *ecdsa.PublicKey) *ecies.PublicKey {
	if pub == nil {
		fmt.Println("ECDSAPublicKeyToECIES: Public key is nil")
		return nil
	}

	xBytes := pub.X.FillBytes(make([]byte, 32))
	yBytes := pub.Y.FillBytes(make([]byte, 32))
	pubBytes := append(xBytes, yBytes...)

	pk, err := ecies.NewPublicKeyFromBytes(pubBytes)
	if err != nil {
		fmt.Println("ECDSAPublicKeyToECIES: Error creating ECIES public key:", err)
		return nil
	}
	return pk
}

func ECDSAPrivateKeyToECIES(priv *ecdsa.PrivateKey) *ecies.PrivateKey {
	if priv == nil {
		return nil
	}
	return ecies.NewPrivateKeyFromBytes(priv.D.Bytes())
}
