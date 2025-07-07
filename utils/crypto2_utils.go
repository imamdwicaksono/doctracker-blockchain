// utils/crypto.go
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math/big"
	"strings"
)

// ECDHKeyPair menghasilkan pasangan kunci ECDH
type ECDHKeyPair struct {
	PrivateKey *big.Int
	PublicKey  *ecdsa.PublicKey
}

// GenerateECDHKeyPair menghasilkan pasangan kunci ECDH
func GenerateECDHKeyPair() (*ECDHKeyPair, error) {
	curve := elliptic.P256()
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECDHKeyPair{
		PrivateKey: new(big.Int).SetBytes(priv),
		PublicKey: &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}, nil
}

// DeriveSharedSecret menghasilkan shared secret menggunakan ECDH
func DeriveSharedSecret(privKey *big.Int, pubKey *ecdsa.PublicKey) ([]byte, error) {
	if pubKey == nil {
		return nil, errors.New("public key is nil")
	}

	x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, privKey.Bytes())
	if x == nil {
		return nil, errors.New("failed to derive shared secret")
	}

	secret := make([]byte, (pubKey.Curve.Params().BitSize+7)/8)
	return x.FillBytes(secret), nil
}

// EncryptData mengenkripsi data dengan AES-GCM
func EncryptData(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// DecryptData mendekripsi data dengan AES-GCM
func DecryptData(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// SerializePublicKey mengubah public key menjadi byte
func SerializePublicKey(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}

	keyLen := (pub.Curve.Params().BitSize + 7) / 8
	serialized := make([]byte, 1+2*keyLen)
	serialized[0] = 0x04 // uncompressed format

	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()

	copy(serialized[1+keyLen-len(xBytes):], xBytes)
	copy(serialized[1+2*keyLen-len(yBytes):], yBytes)

	return serialized
}

// DeserializePublicKey mengubah byte menjadi public key
func DeserializePublicKey(data []byte, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid public key data")
	}

	// Only support uncompressed format
	if data[0] != 0x04 {
		return nil, errors.New("only uncompressed public keys supported")
	}

	keyLen := (curve.Params().BitSize + 7) / 8
	if len(data) != 1+2*keyLen {
		return nil, errors.New("invalid public key length")
	}

	x := new(big.Int).SetBytes(data[1 : 1+keyLen])
	y := new(big.Int).SetBytes(data[1+keyLen:])

	if !curve.IsOnCurve(x, y) {
		return nil, errors.New("point is not on curve")
	}

	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

// GenerateEncryptionKey menghasilkan kunci enkripsi dari shared secret
func GenerateEncryptionKey(sharedSecret []byte) []byte {
	hash := sha256.Sum256(sharedSecret)
	return hash[:]
}

// ECIESPublicKey is a placeholder for ECIES public key type
type ECIESPublicKey ecdsa.PublicKey

// ECIESPrivateKey is a placeholder for ECIES private key type
type ECIESPrivateKey struct {
	D     *big.Int
	Pub   ECIESPublicKey
	Curve elliptic.Curve
}

// ECIESEncrypt mengenkripsi data dengan ECIES
func ECIESEncrypt(rand io.Reader, pub *ECIESPublicKey, data []byte, s1, s2 []byte) ([]byte, error) {
	// Validasi public key
	if pub == nil {
		return nil, errors.New("public key is nil")
	}
	if !pub.Curve.IsOnCurve(pub.X, pub.Y) {
		return nil, errors.New("public key is not on curve")
	}

	// Implementasi enkripsi ECIES sederhana
	// (Dalam produksi, gunakan implementasi lengkap)
	ct := make([]byte, len(data)+32) // +32 untuk IV/key
	if _, err := rand.Read(ct[:32]); err != nil {
		return nil, err
	}

	// Simulasi enkripsi (implementasi aktual lebih kompleks)
	for i := range data {
		ct[i+32] = data[i] ^ ct[i%32]
	}

	return ct, nil
}

// ECIESDecrypt mendekripsi data ECIES
func ECIESDecrypt(priv *ECIESPrivateKey, ct []byte, s1, s2 []byte) ([]byte, error) {
	if priv == nil {
		return nil, errors.New("private key is nil")
	}
	if len(ct) < 32 {
		return nil, errors.New("ciphertext too short")
	}

	// Simulasi dekripsi (implementasi aktual lebih kompleks)
	pt := make([]byte, len(ct)-32)
	for i := range pt {
		pt[i] = ct[i+32] ^ ct[i%32]
	}

	return pt, nil
}

// ImportECDSA converts ECDSA private key to our ECIES format
func ImportECDSA(priv *ecdsa.PrivateKey) *ECIESPrivateKey {
	return &ECIESPrivateKey{
		D: priv.D,
		Pub: ECIESPublicKey{
			X:     priv.PublicKey.X,
			Y:     priv.PublicKey.Y,
			Curve: priv.PublicKey.Curve,
		},
	}
}

func ImportECDSAPublic(pub *ecdsa.PublicKey) *ECIESPublicKey {
	return &ECIESPublicKey{
		X:     pub.X,
		Y:     pub.Y,
		Curve: pub.Curve,
	}
}

// decryptGCM helper for AES-GCM decryption
func decryptGCM(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil, data[:gcm.NonceSize()], data[gcm.NonceSize():], nil)
}

func ConvertForAtoi(input string) (string, error) {
	parts := strings.Split(input, ":")
	if len(parts) < 2 {
		return "", errors.New("invalid format, expected index:hash")
	}

	indexStr := parts[0]
	return indexStr, nil
}
