package models

import "time"

type Tracker struct {
	ID          string `gorm:"primaryKey;size:64" json:"id"`
	Type        string `gorm:"size:50" json:"type"`
	Privacy     string `gorm:"size:20" json:"privacy"`
	Creator     string `gorm:"size:255" json:"creator"`
	CreatorAddr string `gorm:"size:255" json:"creator_address"`

	CreatedAt int64  `json:"created_at"`
	TargetEnd string `gorm:"size:100" json:"target_end"`
	Status    string `gorm:"size:20" json:"status"`

	Checkpoints    []Checkpoint      `gorm:"type:jsonb" json:"checkpoints"`
	EncryptedNotes map[string]string `gorm:"type:jsonb" json:"encrypted_notes"`

	UpdatedAt time.Time `json:"-"`
}

type Checkpoint struct {
	Emails         []string          `json:"emails,omitempty"`  // multiple recipients
	Email          string            `json:"email,omitempty"`   // single recipient (legacy support)
	Type           string            `json:"type"`              // internal / external
	Company        string            `json:"company,omitempty"` // only if external
	Role           string            `json:"role"`              // signer / courier
	IsViewable     bool              `json:"is_view"`           // if true, can decrypt Note
	Note           string            `json:"note,omitempty"`    // original (optional)
	EncryptedNote  string            `json:"encrypted_note,omitempty"`
	EncryptedNotes map[string]string `json:"encrypted_notes,omitempty"` // map[email]encrypted_note
	Address        string            `json:"address,omitempty"`         // first address for backward compat
	Addresses      []string          `json:"addresses,omitempty"`       // all recipient addresses

	EvidenceHash string `json:"evidence_hash,omitempty"`
	EvidencePath string `json:"evidence_path,omitempty"`

	IsCompleted bool  `json:"is_completed"`
	CompletedAt int64 `json:"completed_at,omitempty"`
}

type CheckpointStatusInput struct {
	TrackerID string  `json:"tracker_id"`
	Email     string  `json:"email"`
	Note      *string `json:"note"`
	Evidence  *string `json:"evidence"` // base64 encoded image
}

type CheckpointCompleteRequest struct {
	TrackerID   string  `json:"tracker_id" example:"1234"`
	Email       string  `json:"email" example:"user@example.com"`
	Note        *string `json:"note,omitempty" example:"Dokumen sudah ditandatangani"`
	EvidenceB64 *string `json:"evidence,omitempty" example:"data:image/png;base64,iVBORw0KGgo..."`
}
