package audit

import (
	"crypto/sha256"
	"doc-tracker/blockchain/audit"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func MerkleRoot(hashes []string) string {
	if len(hashes) == 0 {
		return ""
	}

	for len(hashes) > 1 {
		var next []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				h := sha256.Sum256([]byte(hashes[i] + hashes[i+1]))
				next = append(next, hex.EncodeToString(h[:]))
			} else {
				next = append(next, hashes[i])
			}
		}
		hashes = next
	}
	return hashes[0]
}

func BlockHash(height uint64, prev, merkle string, ts int64) string {
	raw := fmt.Sprintf("%d|%s|%s|%d", height, prev, merkle, ts)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func Verify(height uint64, hash string) (bool, error) {
	client, err := audit.NewClient()
	if err != nil {
		return false, err
	}

	return client.Contract.Verify(
		nil,
		big.NewInt(int64(height)),
		common.HexToHash(hash),
	)
}
