package audit

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

type BlockPayload struct {
	Height     uint64   `json:"height"`
	PrevHash   string   `json:"prev_hash"`
	Timestamp  int64    `json:"timestamp"`
	TxHashes   []string `json:"tx_hashes"`
	MerkleRoot string   `json:"merkle_root"`
	BlockHash  string   `json:"block_hash"`
}

func LoadBlock(path string) (*BlockPayload, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var b BlockPayload
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func VerifyLocal(b *BlockPayload) error {
	m := MerkleRoot(b.TxHashes)
	if m != b.MerkleRoot {
		return errors.New("merkle root mismatch")
	}

	h := BlockHash(b.Height, b.PrevHash, m, b.Timestamp)
	if h != b.BlockHash {
		return errors.New("block hash mismatch")
	}
	return nil
}

func StringToHash(s string) common.Hash {
	return common.HexToHash(s)
}
