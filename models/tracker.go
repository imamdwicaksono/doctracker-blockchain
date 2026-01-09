package models

import "time"

type Tracker struct {
	ID             string `gorm:"primaryKey;size:64" json:"id"`
	Type           string `gorm:"size:50" json:"type"`
	Privacy        string `gorm:"size:20" json:"privacy"`
	Creator        string `gorm:"size:255;index" json:"creator"`
	CreatorAddress string `gorm:"size:255" json:"creator_address"`

	TargetEnd string `gorm:"size:100" json:"target_end"`
	Status    string `gorm:"size:20;index" json:"status"`

	CreatedAt int64     `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Checkpoints Checkpoints `gorm:"type:jsonb" json:"checkpoints"`
}

type Checkpoint struct {
	Emails         []string          `json:"emails,omitempty"` // INTERNAL ONLY
	Email          string            `json:"email,omitempty"`  // UTAMA
	Type           string            `json:"type"`
	Company        string            `json:"company,omitempty"`
	Role           string            `json:"role"`
	IsViewable     bool              `json:"is_view"`
	Note           string            `json:"note,omitempty"`
	EncryptedNote  string            `json:"encrypted_note,omitempty"`  // legacy single
	EncryptedNotes map[string]string `json:"encrypted_notes,omitempty"` // INTERNAL
	Address        string            `json:"address,omitempty"`
	Addresses      []string          `json:"addresses,omitempty"` // INTERNAL

	EvidenceHash string `json:"evidence_hash,omitempty"`
	EvidencePath string `json:"evidence_path,omitempty"`

	IsCompleted bool  `json:"is_completed"`
	CompletedAt int64 `json:"completed_at,omitempty"`
}

type CheckpointRecipient struct {
	ID           uint `gorm:"primaryKey"`
	CheckpointID uint `gorm:"index;not null"`

	Email   string `gorm:"index"`
	Address string `gorm:"index"`

	EncryptedNote string

	CreatedAt time.Time
}

type CheckpointStatusInput struct {
	TrackerID string `json:"tracker_id"`
	Email     string `json:"email"`
	Note      string `json:"note"`
	Evidence  string `json:"evidence"` // ✅ bukan pointer
}
