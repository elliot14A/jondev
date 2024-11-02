package models

import (
	"time"

	"github.com/elliot14A/jondev/infrastructure/sqlite/generated"
	"github.com/google/uuid"
)

type Hash struct {
	Value string
	Salt  string
}

type HashStatus struct {
	ID             uuid.UUID
	IsGenerated    bool
	GeneratedAt    *time.Time
	LastVerifiedAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *HashStatus) From(value generated.HashStatus) *HashStatus {
	m.ID = value.ID
	m.IsGenerated = value.IsGenerated
	if value.GeneratedAt.Valid {
		m.GeneratedAt = &value.GeneratedAt.Time
	}
	if value.LastVerifiedAt.Valid {
		m.LastVerifiedAt = &value.LastVerifiedAt.Time
	}
	m.CreatedAt = value.CreatedAt
	m.UpdatedAt = value.UpdatedAt

	return m
}
