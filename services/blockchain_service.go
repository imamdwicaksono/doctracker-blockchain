package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/p2p"
)

func MinePendingTrackers() models.Block {
	pending := mempool.GetAll()
	if len(pending) == 0 {
		return models.Block{}
	}

	// Convert []*models.Tracker to []models.Tracker
	var trackers []models.Tracker
	for _, t := range pending {
		trackers = append(trackers, *t)
	}
	block, _ := blockchain.MineNewBlock(trackers)
	if blockchain.CheckDuplicateBlock(block) {
		return models.Block{} // Block sudah ada, tidak perlu dibuat lagi
	}
	mempool.Clear()

	// TODO: Kirim ke peer (sync P2P)
	p2p.BroadcastNewBlock(block)

	return block
}
