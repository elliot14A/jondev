package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/elliot14A/jondev/application/services"
	"github.com/elliot14A/jondev/domain/pkg"
	"github.com/elliot14A/jondev/infrastructure/hash"
	"github.com/elliot14A/jondev/infrastructure/logger"
	"github.com/elliot14A/jondev/infrastructure/sqlite/actions/hash_status"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateHashCmd)
}

var generateHashCmd = &cobra.Command{
	Use:   "generate-hash",
	Short: "🔐 Generate a secure hash for jondev user access",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateHash(); err != nil {
			log.Fatalf("❌ Failed to generate hash: %v", err)
		}
	},
}

func generateHash() error {
	config, err := pkg.LoadConfig()
	if err != nil {
		return fmt.Errorf("❌ Error loading config: %v", err)
	}

	logger := logger.GetLogger()

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

	hashStatus, err := hashSvc.GetStatus(context.Background())
	if err != nil {
		return fmt.Errorf("❌ Failed to check hash status: %v", err)
	}

	if hashStatus.IsGenerated {
		return fmt.Errorf("⚠️  Hash already exists. To regenerate, you must first manually clear both the database and the hash file")
	}

	logger.Info("🔄 Generating new hash for you...")

	// Generate and store h
	h, err := hashSvc.GenerateHash(context.Background(), config.Hash.Key)
	if err != nil {
		return fmt.Errorf("❌ Error generating hash: %v", err)
	}

	fmt.Printf("✨ Hash generated successfully: %s\n", h.Value)

	if err := hashSvc.StoreHash(context.Background(), h); err != nil {
		return fmt.Errorf("❌ Failed to store hash: %v", err)
	}
	logger.Info("💾 Hash successfully saved to file...")

	// Mark as generated in database
	if err := hashStatusRepo.MarkHashAsGenerated(context.Background()); err != nil {
		os.Remove(config.Hash.FilePath)
		return fmt.Errorf("❌ Failed to mark hash as generated: %v", err)
	}

	// Update last verified timestamp
	if err := hashStatusRepo.UpdateLastVerified(context.Background()); err != nil {
		logger.Warn("⚠️  Warning: Failed to update last verified timestamp: %v", err)
	}

	logger.Info("\n🎉 Successfully generated and stored hash!")
	logger.Info("🔒 Your jondev instance is now secured")
	return nil
}
