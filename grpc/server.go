package grpc

import (
	"context"
	"fmt"
	"net"

	"doc-tracker/mempool"
	pb "doc-tracker/proto"
	"doc-tracker/storage"
	"doc-tracker/utils"

	"google.golang.org/grpc"
	"gorm.io/gorm/clause"
)

type server struct {
	pb.UnimplementedP2PServiceServer
}

func (s *server) SyncTracker(
	ctx context.Context,
	in *pb.Tracker,
) (*pb.Empty, error) {

	tracker := utils.TrackerFromProto(in)

	// DB-first, idempotent
	if err := storage.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).
		Create(&tracker).Error; err != nil {

		return &pb.Empty{}, fmt.Errorf("db error: %w", err)
	}

	// cache ke mempool
	mempool.AddIfNotExists(tracker)

	return &pb.Empty{}, nil
}

func (s *server) SyncMempool(
	ctx context.Context,
	in *pb.TrackerList,
) (*pb.Empty, error) {

	for _, t := range in.Trackers {
		tracker := utils.TrackerFromProto(t)

		_ = storage.DB.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).
			Create(&tracker).Error

		mempool.AddIfNotExists(tracker)
	}

	return &pb.Empty{}, nil
}

func (s *server) Ping(
	ctx context.Context,
	_ *pb.Empty,
) (*pb.Pong, error) {

	return &pb.Pong{Message: "ok"}, nil
}

func StartGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterP2PServiceServer(s, &server{})

	fmt.Println("[gRPC] listening on :" + port)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
