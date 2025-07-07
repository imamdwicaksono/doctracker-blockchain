package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	eciesgo "github.com/ecies/go/v2"
)

func mainold() {
	// Generate 32 random bytes (256-bit key)
	priv, err := eciesgo.GenerateKey()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("TRACKER_SECRET:", priv.Hex())              // => Masukkan ini ke .env
	fmt.Println("Public Key   :", priv.PublicKey.Hex(true)) // => Simpan buat log atau referensi
}

func main() {
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

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
}
