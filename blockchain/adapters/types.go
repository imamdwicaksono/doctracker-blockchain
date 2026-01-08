package adapter

type AuditRecord struct {
	DocumentID string
	Hash       [32]byte
	Action     string
	Actor      string
}
