package hash_status

import (
	"context"
)

func (r *SqliteHashStatusRepository) MarkHashAsGenerated(ctx context.Context) error {
	return r.q.MarkHashAsGenerated(ctx)
}

func (r *SqliteHashStatusRepository) UpdateLastVerified(ctx context.Context) error {
	return r.q.UpdateHashLastVerified(ctx)
}
