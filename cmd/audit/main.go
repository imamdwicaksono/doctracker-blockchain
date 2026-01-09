package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"doc-tracker/blockchain/audit"
	internalAudit "doc-tracker/internal/audit"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	loadEnv()

	switch os.Args[1] {
	case "verify":
		verifyCmd(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: audit-cli <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  verify    Verify a block's audit status")
}

func verifyCmd(args []string) {

	fs := flag.NewFlagSet("verify", flag.ExitOnError)

	blockFile := fs.String("block", "", "path to block json")
	height := fs.Uint64("height", 0, "block height")
	hash := fs.String("hash", "", "block hash")

	fs.Parse(args)

	// MODE FILE
	if *blockFile != "" {
		verifyFromFile(*blockFile)
		return
	}

	// MODE ON-CHAIN
	if *height > 0 && *hash != "" {
		verifyOnChain(*height, *hash)
		return
	}

	log.Fatal("invalid usage: use --block OR (--height AND --hash)")
}

func verifyFromFile(blockFile string) {
	block, err := internalAudit.LoadBlock(blockFile)
	if err != nil {
		log.Fatal(err)
	}

	// 1️⃣ Verify local hash
	if err := internalAudit.VerifyLocal(block); err != nil {
		log.Fatalf("LOCAL VERIFY FAILED: %v", err)
	}

	// 2️⃣ Verify on-chain
	client, err := audit.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ok, err := client.Contract.Verify(
		nil,
		big.NewInt(int64(block.Height)),
		common.HexToHash(block.BlockHash),
	)
	if err != nil {
		log.Fatalf("ON-CHAIN VERIFY ERROR: %v", err)
	}

	if !ok {
		fmt.Println("❌ ON-CHAIN VERIFY FAILED (TAMPERED)")
		return
	}

	fmt.Println("✅ AUDIT VALID (FILE + ON-CHAIN VERIFIED)")
}

func verifyOnChain(height uint64, hash string) {
	client, err := audit.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ok, err := client.Contract.Verify(
		nil,
		big.NewInt(int64(height)),
		common.HexToHash(hash),
	)
	if err != nil {
		log.Fatal("ON-CHAIN VERIFY ERROR:", err)
	}

	if ok {
		fmt.Printf("✅ VERIFIED\nBlock #%d is anchored on-chain\n", height)
	} else {
		fmt.Printf("❌ INVALID\nBlock #%d is NOT anchored\n", height)
	}
}

// ===================== ENV =====================
func loadEnv() {

	wd, _ := os.Getwd()
	envPath := filepath.Join(wd, ".env")

	if _, err := os.Stat(envPath); err == nil {
		_ = godotenv.Overload(envPath)
		fmt.Println("✅ .env loaded")
	}
}
