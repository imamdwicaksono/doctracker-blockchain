package adapter

import (
	"context"
	"database/sql"
	"math/big" // Adjust the import path as necessary

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (l *Listener) SyncFromLastBlock(ctx context.Context, db *sql.DB) error {

	lastBlock := LoadLastBlock(db)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(lastBlock + 1)),
		Addresses: []common.Address{l.Address},
	}

	logs, err := l.Client.Eth.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		l.handleLog(db, vLog)
		SaveLastBlock(db, vLog.BlockNumber)
	}

	return nil
}
