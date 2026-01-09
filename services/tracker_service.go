package services

import (
	"time"

	"doc-tracker/models"
	"doc-tracker/storage"
	"doc-tracker/utils"

	"github.com/gofiber/fiber/v2"
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

func GetTrackers(c *fiber.Ctx) ([]models.Tracker, error) {
	email := c.Query("email")
	return GetTrackersByEmail(email)
}

func GetTrackersByAddress(address string) ([]models.Tracker, error) {
	var trackers []models.Tracker

	err := storage.DB.
		Where("id IN (?)",
			storage.DB.
				Table("checkpoints").
				Select("tracker_id").
				Where("? = ANY (addresses)", address),
		).
		Find(&trackers).Error

	return trackers, err
}

func GetTrackerSummary(email string) (map[string]int64, error) {

	type Result struct {
		Status string
		Count  int64
	}

	var results []Result

	err := storage.DB.Raw(`
		SELECT status, COUNT(DISTINCT id) AS count
		FROM trackers
		WHERE
			creator = @email
			OR EXISTS (
				SELECT 1
				FROM jsonb_array_elements(checkpoints) cp
				WHERE cp->>'email' = @email
				   OR @email = ANY (
						SELECT jsonb_array_elements_text(cp->'emails')
				   )
			)
		GROUP BY status
	`, map[string]interface{}{
		"email": email,
	}).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// default values
	summary := map[string]int64{
		"pending":  0,
		"progress": 0,
		"complete": 0,
	}

	for _, r := range results {
		summary[r.Status] = r.Count
	}

	return summary, nil
}

func GetTrackerByHash(evidenceHash string) (models.Tracker, error) {
	var tracker models.Tracker

	err := storage.DB.
		Joins("JOIN checkpoints c ON c.tracker_id = trackers.id").
		Where("c.evidence_hash = ?", evidenceHash).
		First(&tracker).Error

	if err != nil {
		return models.Tracker{}, utils.ErrNotFound
	}
	return tracker, nil
}

func GetTrackersByEmail(email string) ([]models.Tracker, error) {
	var trackers []models.Tracker

	if err := storage.DB.Order("created_at DESC").Find(&trackers).Error; err != nil {
		return nil, err
	}

	return filterTrackersByEmail(trackers, email), nil
}

func GetDataTrackerFromUserLogin(email string) ([]models.Tracker, error) {
	var all []models.Tracker

	if err := storage.DB.Order("created_at DESC").Find(&all).Error; err != nil {
		return nil, err
	}

	return filterTrackersByEmail(all, email), nil
}

func filterTrackersByEmail(
	trackers []models.Tracker,
	email string,
) []models.Tracker {

	seen := map[string]bool{}
	var result []models.Tracker

	for _, t := range trackers {
		if t.Creator == email {
			seen[t.ID] = true
			result = append(result, t)
			continue
		}

		for _, cp := range t.Checkpoints {
			if cp.Email == email {
				if !seen[t.ID] {
					seen[t.ID] = true
					result = append(result, t)
				}
				break
			}
		}
	}

	return result
}

func GetTrackerByID(id string) (models.Tracker, error) {
	var tracker models.Tracker

	err := storage.DB.
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
