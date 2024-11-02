package repositories

import (
	"context"

	"github.com/elliot14A/jondev/domain/models"
)

type HashRepository interface {
	Store(ctx context.Context, hash models.Hash) error
	Read(ctx context.Context) (models.Hash, error)
}

type HashStatusRepository interface {
	GetHashStatus(ctx context.Context) (*models.HashStatus, error)
	MarkHashAsGenerated(ctx context.Context) error
	UpdateLastVerified(ctx context.Context) error
}
