// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package generated

import (
	"context"
)

type Querier interface {
	// queries/hash_status.sql
	GetHashStatus(ctx context.Context) (HashStatus, error)
	MarkHashAsGenerated(ctx context.Context) error
	UpdateHashLastVerified(ctx context.Context) error
}

var _ Querier = (*Queries)(nil)