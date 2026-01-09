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

		trackers, err := p2p.FetchTrackersFromPeer(peer)
		if err != nil {
			fmt.Printf("[Sync] failed from %s: %v\n", peer, err)
			continue
		}

		for _, t := range trackers {
			// ✅ DB-first, idempotent
			if err := storage.DB.
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					DoNothing: true,
				}).
				Create(&t).Error; err != nil {
				fmt.Printf("[Sync] DB insert failed (%s): %v\n", t.ID, err)
				continue
			}

			// 🟡 mempool hanya cache → best effort
			mempool.AddIfNotExists(t)
		}
	}
}
