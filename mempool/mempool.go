// mempool/mempool.go
package mempool

import (
	"doc-tracker/models"
	"doc-tracker/utils"
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TrackerEntry = models.Tracker

var (
	mempool = make(map[string]*models.Tracker)
	mu      sync.RWMutex
)

// InitKeys inisialisasi atau load kunci
func InitKeys() (*utils.ECDHKeyPair, error) {
	return utils.LoadKeysFromEnv()
}

func InitFromDB(db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	var trackers []models.Tracker
	if err := db.
		Where("status IN ?", []string{"progress", "pending"}).
		Find(&trackers).Error; err != nil {
		return err
	}

	mempool = make(map[string]*models.Tracker)
	for i := range trackers {
		t := trackers[i]
		mempool[t.ID] = &t
	}

	return nil
}

func Add(db *gorm.DB, t *models.Tracker) error {
	mu.Lock()
	defer mu.Unlock()

	if err := db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).
		Create(t).Error; err != nil {
		return err
	}

	mempool[t.ID] = t
	return nil
}

// Get semua tracker di mempool
func GetAll() []*models.Tracker {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]*models.Tracker, 0, len(mempool))
	for _, t := range mempool {
		list = append(list, t)
	}
	return list
}

func GetByID(id string) *models.Tracker {
	mu.RLock()
	defer mu.RUnlock()

	if tracker, exists := mempool[id]; exists {
		return tracker
	}
	return nil
}

// Ambil tracker dengan status complete
func GetCompletedTrackers() []*models.Tracker {
	mu.RLock()
	defer mu.RUnlock()

	var completed []*models.Tracker
	for _, t := range mempool {
		if t.Status == "complete" {
			completed = append(completed, t)
		}
	}
	return completed
}

// Hapus tracker dari mempool
func RemoveFromMempool(db *gorm.DB, id string) {
	Remove(db, id)
}

func SyncMempool(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tracker models.Tracker
		if err := c.BodyParser(&tracker); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid data"})
		}

		if err := Add(db, &tracker); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "ok"})
	}
}

func GetProgressTrackers() []*models.Tracker {
	mu.RLock()
	defer mu.RUnlock()

	var progressTrackers []*models.Tracker
	for _, t := range mempool {
		if t.Status == "progress" {
			progressTrackers = append(progressTrackers, t)
		}
	}
	return progressTrackers
}

func Update(db *gorm.DB, tracker *models.Tracker) error {
	return UpdateTracker(db, tracker)
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	mempool = make(map[string]*models.Tracker)
}

func AddIfNotExists(tracker models.Tracker) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := mempool[tracker.ID]; !exists {
		mempool[tracker.ID] = &tracker
	}
}

func Exists(id string) bool {
	mu.RLock()
	defer mu.RUnlock()

	_, exists := mempool[id]
	return exists
}

// Iterate iterates over all transactions in the mempool and applies the given function.
func Iterate(fn func(tx *models.Tracker) error, email string) error {
	mu.RLock()
	snapshot := make([]*models.Tracker, 0, len(mempool))
	for _, tx := range mempool {
		snapshot = append(snapshot, tx)
	}
	mu.RUnlock()

	for _, tx := range snapshot {
		if err := fn(tx); err != nil {
			return err
		}
	}
	return nil
}

func UpdateTracker(db *gorm.DB, tracker *models.Tracker) error {
	mu.Lock()
	defer mu.Unlock()

	if err := db.Save(tracker).Error; err != nil {
		return err
	}

	mempool[tracker.ID] = tracker
	return nil
}

func Remove(db *gorm.DB, id string) error {
	mu.Lock()
	defer mu.Unlock()

	if err := db.Delete(&models.Tracker{}, "id = ?", id).Error; err != nil {
		return err
	}

	delete(mempool, id)
	return nil
}
