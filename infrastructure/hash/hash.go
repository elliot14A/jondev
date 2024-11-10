package hash

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/elliot14A/jondev/domain/models"
)

type FileHashRepository struct {
	filePath string
	key      []byte
}

func NewHashRepository(filePath string, rawKey string) *FileHashRepository {
	key := sha256.Sum256([]byte(rawKey))
	return &FileHashRepository{filePath: filePath, key: key[:]}
}

func (f *FileHashRepository) Store(ctx context.Context, hash models.Hash) error {
	data := []byte(hash.Value)

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
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return models.Hash{}, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and encrypted data
	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return models.Hash{}, fmt.Errorf("failed to decrypt: %v", err)
	}

	// Parse data using the same separator as Store

	return models.Hash{
		Value: string(plaintext),
	}, nil
}
