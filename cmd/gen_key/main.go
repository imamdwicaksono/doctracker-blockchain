package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	curve := elliptic.P256()

	// Generate private key
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Serialize private key (D)
	privHex := hex.EncodeToString(priv)

	// Serialize public key (uncompressed: 04 || X || Y)
	pub := elliptic.Marshal(curve, x, y)
	pubHex := hex.EncodeToString(pub)

	fmt.Println("ECDH_PRIVATE_KEY=", privHex)
	fmt.Println("ECDH_PUBLIC_KEY=", pubHex)
}
