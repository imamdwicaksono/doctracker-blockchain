package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"

	"github.com/tyler-smith/go-bip39"
)

func PrivateKeyFromEntropy(entropy []byte) *ecdsa.PrivateKey {
	if len(entropy) < 32 {
		return nil
	}
	reader := bytes.NewReader(entropy)
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), reader)
	return priv
}
func RandomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func EncryptWithPublicKey(plainText string, pub *ecdsa.PublicKey) string {
	key := sha256.Sum256(elliptic.Marshal(pub.Curve, pub.X, pub.Y))
	block, _ := aes.NewCipher(key[:])
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainText))
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptWithPrivateKey(cipherTextB64 string, priv *ecdsa.PrivateKey) string {
	cipherText, _ := base64.StdEncoding.DecodeString(cipherTextB64)
	key := sha256.Sum256(elliptic.Marshal(priv.Curve, priv.PublicKey.X, priv.PublicKey.Y))
	block, _ := aes.NewCipher(key[:])
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText)
}

func IsValidMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// PrivateKeyFromPEM parses a PEM encoded private key and returns it as an interface{}.
func PrivateKeyFromPEM(pemStr string) interface{} {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privKey
	}
	// Try PKCS8
	privKey2, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err2 == nil {
		return privKey2
	}
	return nil
}

// DecryptMessage decrypts the given ciphertext using the provided RSA private key.
func DecryptMessage(privKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, ciphertext)
}

func DecryptNoteRSA(encryptedBase64 string, privKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedBase64)
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func EncryptNoteRSA(note string, pubKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pub := pubInterface.(*rsa.PublicKey)

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(note))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func PublicKeyToString(pub *ecdsa.PublicKey) (string, error) {
	if pub == nil {
		return "", errors.New("public key is nil")
	}
	pubBytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	return base64.StdEncoding.EncodeToString(pubBytes), nil
}
