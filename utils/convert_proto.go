package utils

import (
	"doc-tracker/models"
	pb "doc-tracker/proto"
)

func ConvertFromProto(p *pb.Block) models.Block {
	return models.Block{
		Hash:      p.Hash,
		PrevHash:  p.PrevHash,
		Index:     int(p.Index),
		Timestamp: p.Timestamp,
		Nonce:     int(p.Nonce),
		Transactions: func() []models.Tracker {
			txs := make([]models.Tracker, len(p.Transactions))
			for i, tx := range p.Transactions {
				txs[i] = models.Tracker{
					ID:             tx.Id,
					Creator:        tx.Creator,
					Type:           tx.Type,
					Privacy:        tx.Privacy,
					CreatorAddr:    tx.CreatorAddr,
					CreatedAt:      tx.CreatedAt,
					TargetEnd:      tx.TargetEnd,
					Status:         tx.Status,
					EncryptedNotes: tx.EncryptedNotes,
					Checkpoints: func() []models.Checkpoint {
						checkpoints := make([]models.Checkpoint, len(tx.Checkpoints))
						for j, cp := range tx.Checkpoints {
							checkpoints[j] = models.Checkpoint{
								Email:         cp.Email,
								Type:          cp.Type,
								Company:       cp.Company,
								Role:          cp.Role,
								IsViewable:    cp.IsViewable,
								Note:          cp.Note,
								EncryptedNote: cp.EncryptedNote,
								Address:       cp.Address,
								EvidenceHash:  cp.EvidenceHash,
								EvidencePath:  cp.EvidencePath,
								IsCompleted:   cp.IsCompleted,
								CompletedAt:   cp.CompletedAt,
							}
						}
						return checkpoints
					}(),
				}
			}
			return txs
		}(),
	}
}

func ConvertToProto(b models.Block) *pb.Block {
	return &pb.Block{
		Hash:      b.Hash,
		PrevHash:  b.PrevHash,
		Index:     int32(b.Index),
		Timestamp: b.Timestamp,
		Nonce:     int32(b.Nonce),
		Transactions: func() []*pb.Tracker {
			txs := make([]*pb.Tracker, len(b.Transactions))
			for i, tx := range b.Transactions {
				txs[i] = &pb.Tracker{
					Id:             tx.ID,
					Creator:        tx.Creator,
					Type:           tx.Type,
					Privacy:        tx.Privacy,
					CreatorAddr:    tx.CreatorAddr,
					CreatedAt:      tx.CreatedAt,
					TargetEnd:      tx.TargetEnd,
					Status:         tx.Status,
					EncryptedNotes: tx.EncryptedNotes,
					Checkpoints: func() []*pb.Checkpoint {
						checkpoints := make([]*pb.Checkpoint, len(tx.Checkpoints))
						for j, cp := range tx.Checkpoints {
							checkpoints[j] = &pb.Checkpoint{
								Email:         cp.Email,
								Type:          cp.Type,
								Company:       cp.Company,
								Role:          cp.Role,
								IsViewable:    cp.IsViewable,
								Note:          cp.Note,
								EncryptedNote: cp.EncryptedNote,
								Address:       cp.Address,
								EvidenceHash:  cp.EvidenceHash,
								EvidencePath:  cp.EvidencePath,
								IsCompleted:   cp.IsCompleted,
								CompletedAt:   cp.CompletedAt,
							}
						}
						return checkpoints
					}(),
				}
			}
			return txs
		}(),
	}
}

func ConvertToProtoBlockList(blocks []models.Block) *pb.BlockList {
	blockList := &pb.BlockList{
		Blocks: make([]*pb.Block, len(blocks)),
	}

	for i, block := range blocks {
		blockList.Blocks[i] = ConvertToProto(block)
	}

	return blockList
}

func ConvertToProtoBlock(block models.Block) *pb.Block {
	return &pb.Block{
		Hash:      block.Hash,
		PrevHash:  block.PrevHash,
		Index:     int32(block.Index),
		Timestamp: block.Timestamp,
		Nonce:     int32(block.Nonce),
		Transactions: func() []*pb.Tracker {
			txs := make([]*pb.Tracker, len(block.Transactions))
			for i, tx := range block.Transactions {
				txs[i] = &pb.Tracker{
					Id:             tx.ID,
					Creator:        tx.Creator,
					Type:           tx.Type,
					Privacy:        tx.Privacy,
					CreatorAddr:    tx.CreatorAddr,
					CreatedAt:      tx.CreatedAt,
					TargetEnd:      tx.TargetEnd,
					Status:         tx.Status,
					EncryptedNotes: tx.EncryptedNotes,
					Checkpoints:    ConvertToProtoCheckpoints(tx.Checkpoints),
				}
			}
			return txs
		}(),
	}
}

func ConvertToProtoCheckpoints(checkpoints []models.Checkpoint) []*pb.Checkpoint {
	cpList := make([]*pb.Checkpoint, len(checkpoints))
	for i, cp := range checkpoints {
		cpList[i] = &pb.Checkpoint{
			Email:         cp.Email,
			Type:          cp.Type,
			Company:       cp.Company,
			Role:          cp.Role,
			IsViewable:    cp.IsViewable,
			Note:          cp.Note,
			EncryptedNote: cp.EncryptedNote,
			Address:       cp.Address,
			EvidenceHash:  cp.EvidenceHash,
			EvidencePath:  cp.EvidencePath,
			IsCompleted:   cp.IsCompleted,
			CompletedAt:   cp.CompletedAt,
		}
	}
	return cpList
}
