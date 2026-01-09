package p2p

import (
	"context"
	"time"

	pb "doc-tracker/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SyncTrackerToPeer(peerAddr string, tracker *pb.Tracker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial(
		peerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewP2PServiceClient(conn)
	_, err = client.SyncTracker(ctx, tracker)
	return err
}

func SyncMempoolToPeer(peerAddr string, trackers []*pb.Tracker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.Dial(
		peerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewP2PServiceClient(conn)
	_, err = client.SyncMempool(ctx, &pb.TrackerList{
		Trackers: trackers,
	})
	return err
}

func PingPeer(peerAddr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.Dial(
		peerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewP2PServiceClient(conn)
	_, err = client.Ping(ctx, &pb.Empty{})
	return err
}
