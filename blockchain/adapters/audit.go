package adapter

import (
	"context"
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// DocumentAudit ABI binding (hasil abigen)
type DocumentAudit interface {
	Record(
		opts *bind.TransactOpts,
		docId string,
		hash [32]byte,
		action string,
		actor string,
	) (*types.Transaction, error)
}

type AuditService struct {
	Contract DocumentAudit
	Client   *Client
}

func (s *AuditService) RecordDocument(
	ctx context.Context,
	docID string,
	content []byte,
	action string,
	actor string,
) (string, error) {

	hash := sha256.Sum256(content)

	tx, err := s.Contract.Record(
		s.Client.Auth,
		docID,
		hash,
		action,
		actor,
	)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
