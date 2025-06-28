package mempool

import (
	"doc-tracker/models"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber"
)

var mempool = make(map[string]*models.Tracker)

var (
	mempoolList map[string]models.Tracker = make(map[string]models.Tracker)
	mu          sync.Mutex
	filePath    string = "data/mempool.json" // Path to save mempool data
)

type TrackerEntry = models.Tracker

// Add tracker ke mempool jika belum ada
func Add(t *models.Tracker) {
	if _, exists := mempool[t.ID]; !exists {
		mempool[t.ID] = t
	}
	SaveToFile()
}

// Get semua tracker di mempool
func GetAll() []*models.Tracker {
	var list []*models.Tracker
	for _, t := range mempool {
		list = append(list, t)
	}
	return list
}

// Ambil tracker dengan status complete
func GetCompletedTrackers() []*models.Tracker {
	var completed []*models.Tracker
	for _, t := range mempool {
		if t.Status == "complete" {
			completed = append(completed, t)
		}
	}
	return completed
}

// Hapus tracker dari mempool
func RemoveFromMempool(id string) {
	delete(mempool, id)
}

func SyncMempool(c *fiber.Ctx) error {
	var tracker models.Tracker
	if err := c.BodyParser(&tracker); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid tracker data"})
	}

	Add(&tracker) // gunakan fungsi baru
	return c.JSON(fiber.Map{"status": "tracker synced"})
}

func GetProgressTrackers() []*models.Tracker {
	var list []*models.Tracker
	for _, t := range mempool {
		if t.Status == "progress" {
			list = append(list, t)
		}
	}
	return list
}

func Update(tracker models.Tracker) {
	if _, exists := mempool[tracker.ID]; exists {
		mempool[tracker.ID] = &tracker
		return
	}
	// jika tidak ditemukan, tambahkan
	mempool[tracker.ID] = &tracker
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	mempoolList = make(map[string]models.Tracker)
}

func AddIfNotExists(tracker models.Tracker) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := mempoolList[tracker.ID]; !exists {
		mempoolList[tracker.ID] = tracker
	}
}

// Iterate iterates over all transactions in the mempool and applies the given function.
func Iterate(fn func(tx *models.Tracker) error) error {
	for _, tx := range mempool { // assuming mempool is a slice or map of *models.Transaction
		if err := fn(tx); err != nil {
			return err
		}
	}
	return nil
}

func SaveToFile() {
	data, err := json.MarshalIndent(mempool, "", "  ")
	if err != nil {
		log.Println("Failed to marshal mempool:", err)
		return
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Println("Failed to save mempool to file:", err)
	}
}

func LoadFromFile() {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Failed to read mempool file:", err)
		return
	}

	var flat map[string]models.Tracker
	err = json.Unmarshal(data, &flat)
	if err != nil {
		log.Println("Failed to unmarshal mempool:", err)
		return
	}

	for id, val := range flat {
		tracker := val
		mempool[id] = &tracker
	}

	log.Println("Mempool loaded from file")
}
