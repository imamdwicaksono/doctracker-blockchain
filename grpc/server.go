package grpc

import (
	"context"
	"doc-tracker/blockchain"
	pb "doc-tracker/proto" // ganti sesuai path
	"doc-tracker/utils"
	"fmt"

	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedP2PServiceServer
}

func (s *server) BroadcastBlock(ctx context.Context, in *pb.Block) (*pb.Empty, error) {
	block := utils.ConvertFromProto(in) // ✅ convert to internal model
	add := blockchain.TryAddBlock(block)
	if !add {
		return &pb.Empty{}, fmt.Errorf("failed to add block")
	}
	return &pb.Empty{}, nil
}

func (s *server) GetLatestBlock(ctx context.Context, in *pb.Empty) (*pb.Block, error) {
	block := blockchain.GetLastBlock()
	return utils.ConvertToProto(block), nil // ✅ convert back to proto
}

func StartGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterP2PServiceServer(s, &server{})
	s.Serve(lis)
}
