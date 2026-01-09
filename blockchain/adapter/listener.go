package adapter

import (
	"context"
	"doc-tracker/blockchain/audit"
	"doc-tracker/blockchain/audit/contract"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

type Listener struct {
	Client   *audit.Client
	Contract *contract.Contract
	Address  common.Address
}

func (l *Listener) StartPolling(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("📡 Blockchain polling listener started")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Blockchain polling stopped")
			return

		case <-ticker.C:
			if err := l.pollOnce(ctx, db); err != nil {
				log.Println("poll error:", err)
			}
		}
	}
}

func (l *Listener) pollOnce(ctx context.Context, db *gorm.DB) error {

	lastBlock := LoadLastBlock(db)

	header, err := l.Client.Eth.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}

	latest := header.Number.Uint64()
	if lastBlock >= latest {
		return nil
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(lastBlock + 1)),
		ToBlock:   big.NewInt(int64(latest)),
		Addresses: []common.Address{l.Address},
	}

	logs, err := l.Client.Eth.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	var maxBlock uint64 = lastBlock

	for _, vLog := range logs {
		if err := l.handleLog(db, vLog); err != nil {
			log.Println("handle log error:", err)
			continue
		}

		if vLog.BlockNumber > maxBlock {
			maxBlock = vLog.BlockNumber
		}
	}

	// ✅ SAVE ONCE (atomic)
	if maxBlock > lastBlock {
		SaveLastBlock(db, maxBlock)
	}

	return nil
}

func (l *Listener) handleLog(db *gorm.DB, vLog types.Log) error {

	event, err := l.Contract.ParseBlockAnchored(vLog)
	if err != nil {
		return err
	}

	height := event.Height.Uint64()
	hash := common.BytesToHash(event.Hash[:]).Hex()

	// ✅ Idempotent insert

	audit, err := SelectAuditByHeight(db, height)
	if err != nil {
		return err
	}

	if audit == nil {
		SaveAuditIfNotExists(db, height, hash, vLog.TxHash.Hex())
		log.Printf("🔔 Block anchored | height=%d hash=%s", height, hash)
	}

	return err
}

func SaveAuditIfNotExists(
	db *gorm.DB,
	height uint64,
	hash string,
	txHash string,
) error {

	return db.Exec(`
		INSERT INTO blockchain_audits (block_height, block_hash, evm_tx_hash)
		VALUES (?, ?, ?)
		ON CONFLICT (block_height) DO NOTHING
	`, height, hash, txHash).Error
}

func SelectAuditByHeight(
	db *gorm.DB,
	height uint64,
) (*audit.BlockPayload, error) {

	var record struct {
		BlockHeight uint64
		BlockHash   string
	}

	err := db.Raw(`
		SELECT block_height, block_hash
		FROM blockchain_audits
		WHERE block_height = ?
	`, height).Scan(&record).Error
	if err != nil {
		return nil, err
	}
	if record.BlockHeight == 0 {
		return nil, nil
	}

	return &audit.BlockPayload{
		Height:    record.BlockHeight,
		BlockHash: record.BlockHash,
	}, nil
}
