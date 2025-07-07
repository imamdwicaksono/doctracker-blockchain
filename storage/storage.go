package storage

import (
	"doc-tracker/models"
	"encoding/json"
	"os"
	"sync"
)

var trackerFile = "data/trackers.json"

var Trackers = map[string]*models.Tracker{}

var blockStore = struct {
	blocks []models.Block
	sync.Mutex
}{}

func LoadTrackers() error {
	blockStore.Lock()
	defer blockStore.Unlock()

	data, err := os.ReadFile(trackerFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no file yet
		}
		return err
	}

	return json.Unmarshal(data, &Trackers)
}

func FindTrackerByID(id string) *models.Tracker {
	return Trackers[id]
}

func UpdateTracker(t *models.Tracker) {
	Trackers[t.ID] = t
}
