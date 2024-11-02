package hash

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/elliot14A/jondev/application/repositories"
	"github.com/elliot14A/jondev/domain/models"
)

type FileHashRepository struct {
	filePath string
	key      []byte
}

func NewHashRepository(filePath string, key []byte) repositories.HashRepository {
	return &FileHashRepository{filePath: filePath, key: key}
}

func (f *FileHashRepository) Store(ctx context.Context, hash models.Hash) error {
	// Create data string
	data := []byte(fmt.Sprintf("%s\n%s", hash.Value, hash.Salt))

	// Create cipher block
	block, err := aes.NewCipher(f.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %v", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Write to file
	if err := os.WriteFile(f.filePath, ciphertext, 0o600); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func (f *FileHashRepository) Read(ctx context.Context) (models.Hash, error) {
	// Read ciphertext from file
	ciphertext, err := os.ReadFile(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.Hash{}, nil
		}
		return models.Hash{}, fmt.Errorf("failed to read file: %v", err)
	}

	// Create cipher block
	block, err := aes.NewCipher(f.key)
	if err != nil {
		return models.Hash{}, fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return models.Hash{}, fmt.Errorf("failed to create GCM: %v", err)
	}

	// Verify ciphertext length
	if len(ciphertext) < gcm.NonceSize() {
		return models.Hash{}, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce
	nonce, encryptedData := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return models.Hash{}, fmt.Errorf("failed to decrypt: %v", err)
	}

	// Parse data
	parts := strings.Split(string(plaintext), "\n")
	if len(parts) != 2 {
		return models.Hash{}, fmt.Errorf("invalid data format")
	}

	return models.Hash{Value: parts[0], Salt: parts[1]}, nil
}

