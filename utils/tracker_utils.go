package utils

import "doc-tracker/models"

func UniqueTrackers(trackers []models.Tracker) []models.Tracker {
	unique := make([]models.Tracker, 0, len(trackers))
	seen := make(map[string]bool)
	for _, t := range trackers {
		if !seen[t.ID] {
			unique = append(unique, t)
			seen[t.ID] = true
		}
	}
	return unique
}
