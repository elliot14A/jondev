package grpc

import (
	"context"

	"github.com/elliot14A/jondev/application/services"
	"github.com/elliot14A/jondev/infrastructure/logger"
	pb "github.com/elliot14A/jondev/proto/gen/v1/hash"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HashServer struct {
	pb.UnimplementedHashServiceServer
	hashSvc *services.HashService
	logger  *logger.Logger
}

func NewHashServer(svc *services.HashService, logger *logger.Logger) *HashServer {
	return &HashServer{hashSvc: svc, logger: logger}
}

func (s *HashServer) VerifyHash(ctx context.Context, req *pb.VerifyHashRequest) (*pb.VerifyHashResponse, error) {
	if req.Input == "" {
		return nil, status.Error(codes.InvalidArgument, "input cannot be empty")
	}

	storedHash, err := s.hashSvc.ReadHash(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read  stored hash: %v", err)
	}

	matches := s.hashSvc.VerifyHash(ctx, req.Input, storedHash)
	if matches {
		if err := s.hashSvc.UpdateLastVerified(ctx); err != nil {
			s.logger.Info("failed to update last verified timestamp")
		}
	}
	return &pb.VerifyHashResponse{Matches: matches}, nil
}
