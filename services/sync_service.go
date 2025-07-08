package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/p2p"
	"fmt"
	"time"
)

func StartSyncWorker() {
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for {
			<-ticker.C

			fmt.Println("[Sync] Starting initial blockchain sync")
			for _, peer := range p2p.GetPeers() {
				// // Ambil block terbaru
				// latest := p2p.FetchLatestBlockFrom(peer)
				// blockchain.TryAddBlock(latest)

				// // Ambil isi mempool dari peer
				// entries := p2p.FetchMempoolFrom(peer)
				// for _, tx := range entries {
				// 	mempool.AddIfNotExists(tx)
				// }

				mempool.RemoveDuplicateEntries()
				blockchain.RemoveDuplicateBlocks()

				fmt.Printf("[Sync] Fetching block from peer: %s\n", peer)
				block, err := p2p.FetchBlockGRPC(peer)
				if err != nil {
					p2p.WritePeerListToFile("data/peers.txt")
					fmt.Println("Error fetching block:", err)
					continue
				}
				if block != nil {
					// Process the block
					// blockchain.TryAddBlock(block)
				}
			}
		}
	}()
}
