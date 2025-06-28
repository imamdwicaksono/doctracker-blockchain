package models

type Block struct {
	Index        int       `json:"index"`
	Timestamp    int64     `json:"timestamp"`
	PrevHash     string    `json:"prev_hash"`
	Hash         string    `json:"hash"`
	Nonce        int       `json:"nonce"`
	Transactions []Tracker `json:"transactions"`
}
