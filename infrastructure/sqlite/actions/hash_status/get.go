package hash_status

import (
	"context"
	"fmt"

	"github.com/elliot14A/jondev/domain/models"
)

func (r *SqliteHashStatusRepository) GetHashStatus(ctx context.Context) (*models.HashStatus, error) {
	dbHashStatus, err := r.q.GetHashStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting hash status: %v", err)
	}

	var hashStatus models.HashStatus
	h := models.Convert(dbHashStatus, &hashStatus)

	return h, nil
}
