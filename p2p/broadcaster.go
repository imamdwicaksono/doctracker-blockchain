package p2p

import (
	"bytes"
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Broadcast block baru ke semua peer
func BroadcastNewBlock(block blockchain.Block) {
	Broadcast("/p2p/block", block)
}

// Broadcast tracker (mempool entry) ke semua peer
func BroadcastToMempool(entry mempool.TrackerEntry) {
	Broadcast("/p2p/mempool", entry)
}

func Broadcast(path string, data interface{}) {
	for _, peer := range Peers {
		go func(p string) {
			url := fmt.Sprintf("http://%s%s", p, path)
			jsonData, _ := json.Marshal(data)
			http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		}(peer)
	}
}

func FetchLatestBlockFrom(peer string) models.Block {
	resp, err := http.Get(fmt.Sprintf("http://%s/p2p/latest-block", peer))
	if err != nil {
		return models.Block{}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var block models.Block
	json.Unmarshal(body, &block)
	return block
}

func FetchMempoolFrom(peer string) []mempool.TrackerEntry {
	resp, err := http.Get(fmt.Sprintf("http://%s/p2p/mempool", peer))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var entries []mempool.TrackerEntry
	json.Unmarshal(body, &entries)
	return entries
}

func ReceiveMempool(c *fiber.Ctx) error {
	var tx mempool.TrackerEntry
	if err := c.BodyParser(&tx); err != nil {
		return c.SendStatus(400)
	}

	mempool.AddIfNotExists(tx)
	return c.SendStatus(200)
}
