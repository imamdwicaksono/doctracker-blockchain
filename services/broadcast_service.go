// services/broadcast_service.go
package services

import (
	"bytes"
	"doc-tracker/models"
	"encoding/json"
	"fmt"
	"net/http"
)

var KnownPeers = []string{
	"http://localhost:3001",
	"http://localhost:3002",
	// tambahkan node lain di jaringan kamu
}

func BroadcastToMempool(tracker *models.Tracker) {
	body, _ := json.Marshal(tracker)

	for _, peer := range KnownPeers {
		go func(peer string) {
			url := peer + "/api/sync/mempool"
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Failed to send tracker to peer:", peer, err)
				return
			}
			defer resp.Body.Close()
			fmt.Println("Sent tracker to:", peer, resp.Status)
		}(peer)
	}
}
