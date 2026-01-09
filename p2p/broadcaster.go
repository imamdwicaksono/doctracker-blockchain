package p2p

import (
	"context"
	"time"

	"doc-tracker/models"
	"doc-tracker/proto"
	"doc-tracker/storage"
	"doc-tracker/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm/clause"
)

func FetchTrackersFromPeer(peerAddr string) ([]models.Tracker, error) {

	conn, err := grpc.Dial(
		peerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := proto.NewP2PServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 🔒 DATA ONLY (no mempool wording)
	res, err := client.GetTrackers(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}

	var trackers []models.Tracker

	for _, t := range res.Trackers {
		tracker := utils.FromProtoTracker(t)

		// ✅ DB-first, idempotent
		if err := storage.DB.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).
			Create(&tracker).Error; err != nil {
			return nil, err
		}

		trackers = append(trackers, tracker)
	}

	return trackers, nil
}
