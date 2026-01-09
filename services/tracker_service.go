package services

import (
	"time"

	"doc-tracker/models"
	"doc-tracker/storage"
	"doc-tracker/utils"

	"github.com/google/uuid"
)

func CreateTracker(input models.Tracker) (models.Tracker, error) {

	input.ID = uuid.NewString()
	input.CreatedAt = time.Now().Unix()
	input.Status = "progress"

	// Enkripsi checkpoint notes
	for i, cp := range input.Checkpoints {

		encryptedNotes := make(map[string]string)
		var addresses []string

		for _, email := range cp.Emails {
			if email == "" {
				continue
			}

			wallet := GetOrCreateWallet(email)
			addresses = append(addresses, wallet.Address)

			if cp.IsViewable && cp.Note != "" {
				encryptedNotes[email] =
					utils.EncryptWithPublicKey(cp.Note, wallet.PublicKey)
			}
		}

		input.Checkpoints[i].Addresses = addresses
		input.Checkpoints[i].EncryptedNotes = encryptedNotes
		input.Checkpoints[i].Note = ""
		if len(addresses) > 0 {
			input.Checkpoints[i].Address = addresses[0]
		}
	}

	if err := storage.DB.Create(&input).Error; err != nil {
		return models.Tracker{}, err
	}

	return input, nil
}

func GetTrackersByEmail(email string) ([]models.Tracker, error) {
	var trackers []models.Tracker

	err := storage.DB.
		Preload("Checkpoints").
		Where("creator = ?", email).
		Or("id IN (?)",
			storage.DB.
				Table("checkpoints").
				Select("tracker_id").
				Where("email = ?", email),
		).
		Find(&trackers).Error

	return trackers, err
}

func GetTrackerByID(id string) (models.Tracker, error) {
	var tracker models.Tracker

	err := storage.DB.
		Preload("Checkpoints").
		Where("id = ?", id).
		First(&tracker).Error

	if err != nil {
		return models.Tracker{}, utils.ErrNotFound
	}
	return tracker, nil
}

func GetTrackerByEvidenceHash(hash string) (models.Tracker, error) {
	var tracker models.Tracker

	err := storage.DB.
		Joins("JOIN checkpoints c ON c.tracker_id = trackers.id").
		Where("c.evidence_hash = ?", hash).
		First(&tracker).Error

	if err != nil {
		return models.Tracker{}, utils.ErrNotFound
	}
	return tracker, nil
}
