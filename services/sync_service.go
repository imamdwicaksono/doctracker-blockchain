package services

import (
	"fmt"
	"time"

	"doc-tracker/mempool"
	"doc-tracker/p2p"
	"doc-tracker/storage"

	"gorm.io/gorm/clause"
)

func StartSyncWorker() {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for range ticker.C {
			syncFromPeers()
		}
	}()
}

func syncFromPeers() {
	peers := p2p.GetPeers()

	for _, peer := range peers {
		fmt.Printf("[Sync] syncing trackers from peer: %s\n", peer)

		trackers := p2p.FetchTrackersGRPC(peer)

		for _, t := range trackers {
			// DB-first, idempotent
			_ = storage.DB.
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					DoNothing: true,
				}).
				Create(&t).Error

			// cache ke mempool (opsional)
			mempool.AddIfNotExists(t)
		}
	}
}
