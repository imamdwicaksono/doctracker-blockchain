package adapter

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

type Listener struct {
	Client   *Client
	Contract *DocumentAudit
	Address  common.Address
}

func (l *Listener) StartPolling(ctx context.Context, db *gorm.DB) {

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	log.Println("📡 Blockchain polling listener started")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Blockchain polling stopped")
			return

		case <-ticker.C:
			l.pollOnce(ctx, db)
		}
	}
}

func (l *Listener) pollOnce(ctx context.Context, db *gorm.DB) {

	lastBlock := LoadLastBlock(db)

	header, err := l.Client.Eth.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Println("header error:", err)
		return
	}

	latest := header.Number.Uint64()
	if lastBlock >= latest {
		return
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(lastBlock + 1)),
		ToBlock:   big.NewInt(int64(latest)),
		Addresses: []common.Address{l.Address},
	}

	logs, err := l.Client.Eth.FilterLogs(ctx, query)
	if err != nil {
		log.Println("filter logs error:", err)
		return
	}

	for _, vLog := range logs {
		l.handleLog(db, vLog)
		SaveLastBlock(db, vLog.BlockNumber)
	}
}

func (l *Listener) handleLog(db *gorm.DB, vLog types.Log) {
	log.Printf("🔔 New event log: %+v\n", vLog)
}
