package utils

import (
	"doc-tracker/models"
	pb "doc-tracker/proto"
)

/* ===================== TRACKER ===================== */

func ToProtoTracker(t models.Tracker) *pb.Tracker {
	return &pb.Tracker{
		Id:          t.ID,
		Type:        t.Type,
		Privacy:     t.Privacy,
		Creator:     t.Creator,
		CreatorAddr: t.CreatorAddress,
		CreatedAt:   t.CreatedAt,
		TargetEnd:   t.TargetEnd,
		Status:      t.Status,
		Checkpoints: ToProtoCheckpoints(t.Checkpoints),
	}
}

func FromProtoTracker(p *pb.Tracker) models.Tracker {
	return models.Tracker{
		ID:             p.Id,
		Type:           p.Type,
		Privacy:        p.Privacy,
		Creator:        p.Creator,
		CreatorAddress: p.CreatorAddr,
		CreatedAt:      p.CreatedAt,
		TargetEnd:      p.TargetEnd,
		Status:         p.Status,
		Checkpoints:    FromProtoCheckpoints(p.Checkpoints),
	}
}

/* ===================== CHECKPOINT ===================== */

func ToProtoCheckpoints(cps []models.Checkpoint) []*pb.Checkpoint {
	list := make([]*pb.Checkpoint, 0, len(cps))

	for _, cp := range cps {
		list = append(list, &pb.Checkpoint{
			Email:         cp.Email, // ✅ PRIMARY
			Type:          cp.Type,
			Company:       cp.Company,
			Role:          cp.Role,
			IsViewable:    cp.IsViewable,
			Note:          cp.Note,
			EncryptedNote: cp.EncryptedNote, // ✅ legacy
			Address:       cp.Address,
			EvidenceHash:  cp.EvidenceHash,
			EvidencePath:  cp.EvidencePath,
			IsCompleted:   cp.IsCompleted,
			CompletedAt:   cp.CompletedAt,
		})
	}

	return list
}

func FromProtoCheckpoints(cps []*pb.Checkpoint) []models.Checkpoint {
	list := make([]models.Checkpoint, 0, len(cps))

	for _, cp := range cps {
		list = append(list, models.Checkpoint{
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
		})
	}

	return list
}
