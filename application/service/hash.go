package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/elliot14A/jondev/application/repositories"
	"github.com/elliot14A/jondev/domain/models"
	"github.com/xdg-go/pbkdf2"
)

type HashService struct {
	fileStore repositories.HashRepository
	dbStore   repositories.HashStatusRepository
}

func NewHashService(fileStore repositories.HashRepository, dbStore repositories.HashStatusRepository) *HashService {
	return &HashService{fileStore, dbStore}
}

func (s *HashService) GenerateHash(ctx context.Context, input string) (models.Hash, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return models.Hash{}, err
	}

	saltStr := base64.StdEncoding.EncodeToString(salt)
	hash := pbkdf2.Key([]byte(input), salt, 480000, 32, sha256.New)
	hashStr := base64.StdEncoding.EncodeToString(hash)

	return models.Hash{Value: hashStr, Salt: saltStr}, nil
}

func (s *HashService) VerifyHash(ctx context.Context, input string, storedHash models.Hash) bool {
	salt, err := base64.StdEncoding.DecodeString(storedHash.Salt)
	if err != nil {
		return false
	}

	hash := pbkdf2.Key([]byte(input), salt, 480000, 32, sha256.New)
	computedHash := base64.StdEncoding.EncodeToString(hash)

	return computedHash == storedHash.Value
}

func (s *HashService) StoreHash(ctx context.Context, hash models.Hash) error {
	return s.fileStore.Store(ctx, hash)
}

func (s *HashService) ReadHash(ctx context.Context) (models.Hash, error) {
	return s.fileStore.Read(ctx)
}

func (s *HashService) GetStatus(ctx context.Context) (*models.HashStatus, error) {
	return s.dbStore.GetHashStatus(ctx)
}

func (s *HashService) MarkHashAsGenerated(ctx context.Context) error {
	return s.dbStore.MarkHashAsGenerated(ctx)
}

func (s *HashService) UpdateLastVerified(ctx context.Context) error {
	return s.dbStore.UpdateLastVerified(ctx)
}
