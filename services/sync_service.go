package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/p2p"
	"time"
)

func StartSyncWorker() {
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for {
			<-ticker.C

			for _, peer := range p2p.GetPeers() {
				// Ambil block terbaru
				latest := p2p.FetchLatestBlockFrom(peer)
				blockchain.TryAddBlock(latest)

				// Ambil isi mempool dari peer
				entries := p2p.FetchMempoolFrom(peer)
				for _, tx := range entries {
					mempool.AddIfNotExists(tx)
				}
			}
		}
	}()
}
