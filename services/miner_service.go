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

			trackerIdList := make([]string, len(trackerList))
			// Mine block
			trackers := make([]models.Tracker, len(trackerList))
			for i, t := range trackerList {
				trackers[i] = *t
			}
			// newBlock := blockchain.NewBlockFromTransactions(trackers)

			for _, t := range trackerList {
				trackerIdList = append(trackerIdList, t.ID)
			}

			if CheckDuplicateTracker(trackerIdList) {
				fmt.Println("Duplicate tracker found, skipping mining")
				continue
			}

			fmt.Printf("Mining new block with transactions: %v\n", trackerIdList)
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

func CheckDuplicateTracker(trackerIdList []string) bool {
	for _, id := range trackerIdList {
		if blockchain.IsTrackerInBlockchain(id) {
			return true // Tracker sudah ada di blockchain
		}
	}
	return false // Tidak ada duplikat di blockchain
}
