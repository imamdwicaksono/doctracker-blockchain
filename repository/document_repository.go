package repository

import (
	"database/sql"
	"doc-tracker/types"
)

func CreateDocument(db *sql.DB, doc types.Document) error {
	_, err := db.Exec(`
	  INSERT INTO documents (id, title, content_hash, status, created_at)
	  VALUES ($1,$2,$3,$4,now())
	`, doc.ID, doc.Title, doc.Hash, doc.Status)
	return err
}

func UpdateStatus(db *sql.DB, id string, status string) error {
	_, err := db.Exec(`
	  UPDATE documents SET status=$1, updated_at=now()
	  WHERE id=$2
	`, status, id)
	return err
}

func InsertLog(db *sql.DB, docID string, action string, actor string, txHash string) error {
	_, err := db.Exec(`
	  INSERT INTO document_logs (doc_id, action, actor, tx_hash, created_at)
	  VALUES ($1,$2,$3,$4,now())
	`, docID, action, actor, txHash)
	return err
}
