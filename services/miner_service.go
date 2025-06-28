package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/p2p"
	"fmt"
	"time"
)

func StartMinerWorker() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			trackerList := mempool.GetCompletedTrackers()
			if len(trackerList) == 0 {
				continue
			}

			// Mine block
			trackers := make([]models.Tracker, len(trackerList))
			for i, t := range trackerList {
				trackers[i] = *t
			}
			newBlock := blockchain.NewBlockFromTransactions(trackers)

			// Tambahkan ke chain lokal
			blockchain.AddBlock(newBlock)

			fmt.Printf("✅ Block mined: #%d | Hash: %s | Tx: %d\n", newBlock.Index, newBlock.Hash, len(newBlock.Transactions))

			// Hapus tracker dari mempool
			for _, t := range trackerList {
				mempool.RemoveFromMempool(t.ID)
			}

			// Broadcast block ke semua peer
			go p2p.BroadcastNewBlock(newBlock)
		}
	}()
}
