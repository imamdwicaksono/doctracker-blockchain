package utils

import (
	"doc-tracker/models"
	pb "doc-tracker/proto"
)

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

func TrackerFromProto(p *pb.Tracker) models.Tracker {
	return models.Tracker{
		ID:             p.Id,
		Creator:        p.Creator,
		Type:           p.Type,
		Privacy:        p.Privacy,
		CreatorAddr:    p.CreatorAddr,
		CreatedAt:      p.CreatedAt,
		TargetEnd:      p.TargetEnd,
		Status:         p.Status,
		EncryptedNotes: p.EncryptedNotes,
		Checkpoints: func() []models.Checkpoint {
			cps := make([]models.Checkpoint, len(p.Checkpoints))
			for i, cp := range p.Checkpoints {
				cps[i] = models.Checkpoint{
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
			return cps
		}(),
	}
}

func TrackerToProto(t models.Tracker) *pb.Tracker {
	return &pb.Tracker{
		Id:             t.ID,
		Creator:        t.Creator,
		Type:           t.Type,
		Privacy:        t.Privacy,
		CreatorAddr:    t.CreatorAddr,
		CreatedAt:      t.CreatedAt,
		TargetEnd:      t.TargetEnd,
		Status:         t.Status,
		EncryptedNotes: t.EncryptedNotes,
		Checkpoints: func() []*pb.Checkpoint {
			cps := make([]*pb.Checkpoint, len(t.Checkpoints))
			for i, cp := range t.Checkpoints {
				cps[i] = &pb.Checkpoint{
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
			return cps
		}(),
	}
}
