package grpc

import (
	"context"
	"fmt"

	"github.com/elliot14A/jondev/application/services"
	pb "github.com/elliot14A/jondev/proto/gen/v1/hash"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HashServer struct {
	pb.UnimplementedHashServiceServer
	hashSvc *services.HashService
}

func NewHashServer(svc *services.HashService) *HashServer {
	return &HashServer{hashSvc: svc}
}

func (s *HashServer) VerifyHash(ctx context.Context, req *pb.VerifyHashRequest) (*pb.VerifyHashResponse, error) {
	if req.Input == "" {
		return nil, status.Error(codes.InvalidArgument, "input cannot be empty")
	}

	fmt.Println("===>", req.Input)

	storedHash, err := s.hashSvc.ReadHash(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read  stored hash: %v", err)
	}

	fmt.Println("==>", storedHash.Value)

	matches := s.hashSvc.VerifyHash(ctx, req.Input, storedHash)
	if matches {
		if err := s.hashSvc.UpdateLastVerified(ctx); err != nil {
			fmt.Println("failed to update last verified timestamp")
		}
	}
	return &pb.VerifyHashResponse{Matches: matches}, nil
}
