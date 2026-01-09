package blockchain

import (
	"crypto/sha256"
	"doc-tracker/blockchain/audit"
	"doc-tracker/models"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Block = models.Block

func HashTracker(t models.Tracker) string {
	payload := struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Creator   string `json:"creator"`
		Status    string `json:"status"`
		CreatedAt int64  `json:"created_at"`
	}{
		ID:        t.ID,
		Type:      t.Type,
		Creator:   t.Creator,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
	}

	data, _ := json.Marshal(payload)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func MerkleRoot(input []string) string {
	if len(input) == 0 {
		return ""
	}

	hashes := append([]string{}, input...) // COPY

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

func BlockHash(
	height uint64,
	prevHash string,
	merkleRoot string,
	timestamp int64,
) string {
	raw := fmt.Sprintf("%d|%s|%s|%d",
		height,
		prevHash,
		merkleRoot,
		timestamp,
	)

	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func AnchorBlockToChain(
	height uint64,
	blockHash string,
) error {

	client, err := audit.NewClient()
	if err != nil {
		return err
	}

	return client.Anchor(height, blockHash)
}
