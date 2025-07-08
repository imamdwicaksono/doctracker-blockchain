package p2p

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "doc-tracker/proto" // ganti sesuai path
)

func BroadcastToPeer(peerAddr string, entry *pb.Block) error {
	conn, err := grpc.Dial(peerAddr, grpc.WithInsecure()) // gunakan TLS untuk prod
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewP2PServiceClient(conn)
	_, err = client.BroadcastBlock(context.Background(), entry)
	return err
}

func FetchBlockGRPC(peer string) (*pb.BlockList, error) {
	conn, err := grpc.Dial(peer, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect to peer:", peer)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewP2PServiceClient(conn)
	res, err := client.GetBlockchain(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println("Error fetching blockchain:", err)
		return nil, err
	}

	return res, nil
}
