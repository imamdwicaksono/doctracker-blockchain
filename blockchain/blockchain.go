package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"doc-tracker/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitChainFromDB(db *gorm.DB) {
	var count int64
	db.Model(&models.Block{}).Count(&count)

	if count == 0 {
		genesis := models.Block{
			Height:    0,
			BlockHash: hashSHA256("GENESIS"),
			PrevHash:  "",
			TxCount:   0,
			CreatedAt: time.Now().Unix(),
		}
		db.Create(&genesis)
	}
}

func CreateBlock(
	db *gorm.DB,
	trackers []models.Tracker,
) (*models.Block, error) {

	var block models.Block

	err := db.Transaction(func(tx *gorm.DB) error {

		var last models.Block
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Order("height desc").
			Limit(1).
			Take(&last).Error; err != nil {
			return err
		}

		createdAt := time.Now().Unix()
		height := last.Height + 1

		hash := calculateBlockHash(
			last.BlockHash,
			int(height),
			len(trackers),
			createdAt,
		)

		block = models.Block{
			Height:    height,
			BlockHash: hash,
			PrevHash:  last.BlockHash,
			TxCount:   len(trackers),
			CreatedAt: createdAt,
		}

		return tx.Create(&block).Error
	})

	return &block, err
}

func calculateBlockHash(
	prevHash string,
	height int,
	txCount int,
	createdAt int64,
) string {

	data := fmt.Sprintf(
		"%s|%d|%d|%d",
		prevHash,
		height,
		txCount,
		createdAt,
	)

	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func SaveAudit(
	db *gorm.DB,
	block models.Block,
	evmTxHash string,
) error {

	return db.Create(&models.BlockchainAudit{
		BlockHeight: int(block.Height),
		BlockHash:   block.BlockHash,
		EvmTxHash:   evmTxHash,
	}).Error
}

func hashSHA256(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

func GetLatestBlock(db *gorm.DB) (models.Block, error) {
	var block models.Block

	err := db.
		Order("height desc").
		Limit(1).
		Take(&block).Error

	return block, err
}

func GetLastBlock(db *gorm.DB) (models.Block, error) {
	var block models.Block

	err := db.
		Order("height desc").
		Limit(1).
		Take(&block).Error

	return block, err
}

func AddBlock(db *gorm.DB, block *models.Block) error {
	return db.Create(block).Error
}

func IsBlockValid(newBlock, oldBlock models.Block) bool {
	if oldBlock.Height+1 != newBlock.Height {
		return false
	}

	if oldBlock.BlockHash != newBlock.PrevHash {
		return false
	}

	expectedHash := calculateBlockHash(
		newBlock.PrevHash,
		int(newBlock.Height),
		newBlock.TxCount,
		newBlock.CreatedAt,
	)

	if expectedHash != newBlock.BlockHash {
		return false
	}

	return true
}
