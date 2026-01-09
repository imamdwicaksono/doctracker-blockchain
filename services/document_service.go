package services

import (
	"context"
	"database/sql"
	"log"

	"doc-tracker/blockchain/adapter"
	"doc-tracker/repository"
)

func ApproveDocument(
	db *sql.DB,
	chain *adapter.AuditService,
	docID string,
	content []byte,
	actor string,
) error {

	// 1️⃣ update DB dulu
	if err := repository.UpdateStatus(db, docID, "APPROVED"); err != nil {
		return err
	}

	// 2️⃣ emit blockchain log
	txHash, err := chain.RecordDocument(
		context.Background(),
		docID,
		content,
		"APPROVED",
		actor,
	)
	if err != nil {
		log.Println("⚠️ blockchain failed, will retry:", err)
	}

	// 3️⃣ simpan tx hash (opsional)
	repository.InsertLog(db, docID, "APPROVED", actor, txHash)

	return nil
}
