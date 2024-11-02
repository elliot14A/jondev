// Update file: cmd/hash.go

package cmd

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/elliot14A/jondev/application/service"
	"github.com/elliot14A/jondev/domain/pkg"
	"github.com/elliot14A/jondev/infrastructure/hash"
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

	// Ensure key is exactly 32 bytes (256 bits) for AES-256
	key := sha256.Sum256([]byte(config.Hash.Key))

	// Ensure directory exists
	dir := filepath.Dir(config.Hash.FilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("❌ Failed to create directory: %v", err)
	}

	hashRepository := hash.NewHashRepository(config.Hash.FilePath, key[:])

	db, err := sql.Open("sqlite3", "./.jondev/sqlite/jondev.db")
	if err != nil {
		return fmt.Errorf("❌ Failed to open database: %v", err)
	}

	hashStatusRepo := hash_status.NewHashStatusRepository(db)
	hashSvc := service.NewHashService(hashRepository, hashStatusRepo)

	hashStatus, err := hashSvc.GetStatus(context.Background())
	if err != nil {
		return fmt.Errorf("❌ Failed to check hash status: %v", err)
	}

	if hashStatus.IsGenerated {
		return fmt.Errorf("⚠️  Hash already exists. To regenerate, you must first manually clear both the database and the hash file")
	}

	fmt.Println("🔄 Generating new hash for you...")

	// Generate and store hash
	hash, err := hashSvc.GenerateHash(context.Background(), "secure hash")
	if err != nil {
		return fmt.Errorf("❌ Error generating hash: %v", err)
	}

	fmt.Println("✨ Hash generated successfully...")

	if err := hashSvc.StoreHash(context.Background(), hash); err != nil {
		return fmt.Errorf("❌ Failed to store hash: %v", err)
	}

	fmt.Println("💾 Hash successfully saved to file...")

	// Verify storage
	storedHash, err := hashSvc.ReadHash(context.Background())
	if err != nil {
		return fmt.Errorf("❌ Error reading stored hash: %v", err)
	}

	if storedHash.Value != hash.Value || storedHash.Salt != hash.Salt {
		os.Remove(config.Hash.FilePath)
		return fmt.Errorf("❌ Hash verification failed: stored hash doesn't match generated hash")
	}

	fmt.Println("✅ Hash verified successfully...")

	// Mark as generated in database
	if err := hashStatusRepo.MarkHashAsGenerated(context.Background()); err != nil {
		os.Remove(config.Hash.FilePath)
		return fmt.Errorf("❌ Failed to mark hash as generated: %v", err)
	}

	// Update last verified timestamp
	if err := hashStatusRepo.UpdateLastVerified(context.Background()); err != nil {
		log.Printf("⚠️  Warning: Failed to update last verified timestamp: %v", err)
	}

	fmt.Println("\n🎉 Successfully generated and stored hash!")
	fmt.Println("🔒 Your jondev instance is now secured")
	return nil
}
