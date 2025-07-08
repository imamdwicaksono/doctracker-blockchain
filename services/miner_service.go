package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/p2p"
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
			// newBlock := blockchain.NewBlockFromTransactions(trackers)

			// Filter duplikat tracker
			if FilterDuplicateTrackers(trackers) == nil {
				continue
			}
			// Tambahkan ke chain lokal
			mine, error := blockchain.MineNewBlock(trackers)

			if error != nil {
				// Broadcast block ke semua peer
				go p2p.BroadcastNewBlock(mine)
			}

		}
	}()
}

func FilterDuplicateTrackers(trackers []models.Tracker) []models.Tracker {
	unique := make(map[string]bool)
	var filtered []models.Tracker

	for _, t := range trackers {
		if !unique[t.ID] {
			unique[t.ID] = true
			filtered = append(filtered, t)
		}
	}
	return filtered
}
