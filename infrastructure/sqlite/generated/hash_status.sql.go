// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: hash_status.sql

package generated

import (
	"context"
)

const getHashStatus = `-- name: GetHashStatus :one

select id, is_generated, generated_at, last_verified_at, created_at, updated_at from hash_status
limit 1
`

// queries/hash_status.sql
func (q *Queries) GetHashStatus(ctx context.Context) (HashStatus, error) {
	row := q.db.QueryRowContext(ctx, getHashStatus)
	var i HashStatus
	err := row.Scan(
		&i.ID,
		&i.IsGenerated,
		&i.GeneratedAt,
		&i.LastVerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const markHashAsGenerated = `-- name: MarkHashAsGenerated :exec
update hash_status
set is_generated = true,
    generated_at = current_timestamp
where is_generated = false
`

func (q *Queries) MarkHashAsGenerated(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, markHashAsGenerated)
	return err
}

const updateHashLastVerified = `-- name: UpdateHashLastVerified :exec
update hash_status
set last_verified_at = current_timestamp
`

func (q *Queries) UpdateHashLastVerified(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, updateHashLastVerified)
	return err
}
