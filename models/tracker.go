package models

type Tracker struct {
	ID             string            `json:"id"`
	Type           string            `json:"type"`
	Privacy        string            `json:"privacy"`
	Creator        string            `json:"creator"`
	CreatorAddr    string            `json:"creator_address"`
	CreatedAt      int64             `json:"created_at"`
	Checkpoints    []Checkpoint      `json:"checkpoints"`
	TargetEnd      string            `json:"target_end"` // self / email / address
	Status         string            `json:"status"`     // pending, progress, complete
	EncryptedNotes map[string]string `json:"encrypted_notes,omitempty"`
}

type Checkpoint struct {
	Email         string `json:"email"`
	Type          string `json:"type"`    // internal / external
	Company       string `json:"company"` // only if external
	Role          string `json:"role"`    // signer / courier
	IsViewable    bool   `json:"is_view"` // if true, can decrypt Note
	Note          string `json:"note"`    // original (optional)
	EncryptedNote string `json:"encrypted_note"`
	Address       string `json:"address"` // auto-generated

	EvidenceHash string `json:"evidence_hash,omitempty"`
	EvidencePath string `json:"evidence_path,omitempty"`

	IsCompleted bool  `json:"is_completed"`
	CompletedAt int64 `json:"completed_at,omitempty"`
}

type CheckpointStatusInput struct {
	TrackerID   string `json:"tracker_id"`
	Email       string `json:"email"`
	Note        string `json:"note"`
	EvidenceB64 string `json:"evidence"` // base64 encoded image
}
