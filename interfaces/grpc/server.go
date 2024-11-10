package grpc

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/elliot14A/jondev/application/services"
	"github.com/elliot14A/jondev/domain/pkg"
	"github.com/elliot14A/jondev/infrastructure/hash"
	"github.com/elliot14A/jondev/infrastructure/sqlite/actions/hash_status"
	pb "github.com/elliot14A/jondev/proto/gen/v1/hash"
	"google.golang.org/grpc"
)

func RunGrpcServer(config pkg.Config) error {
	serverAddr := config.GetServerAddr()

	// Ensure directory exists
	dir := filepath.Dir(config.Hash.FilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("❌ Failed to create directory: %v", err)
	}

	hashRepository := hash.NewHashRepository(config.Hash.FilePath, config.Hash.Key)
	db, err := sql.Open("sqlite3", "./.jondev/sqlite/jondev.db")
	if err != nil {
		return fmt.Errorf("❌ Failed to open database: %v", err)
	}
	hashStatusRepo := hash_status.NewHashStatusRepository(db)
	hashSvc := services.NewHashService(hashRepository, hashStatusRepo)

	grpcServer := grpc.NewServer()
	hashServer := NewHashServer(hashSvc)

	pb.RegisterHashServiceServer(grpcServer, hashServer)

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	log.Printf("🚀 Starting gRPC server on %s\n", serverAddr)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("❌ Failed to serve: %v", err)
	}

	return nil
}
