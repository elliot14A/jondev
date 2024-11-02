package hash_status

import (
	"database/sql"

	"github.com/elliot14A/jondev/application/repositories"
	"github.com/elliot14A/jondev/infrastructure/sqlite/generated"
)

var defaultStatusID = "67b4ed39-b6b9-4957-9e3a-0938f2ac0ebd"

type SqliteHashStatusRepository struct {
	q *generated.Queries
}

func NewHashStatusRepository(db *sql.DB) repositories.HashStatusRepository {
	return &SqliteHashStatusRepository{
		q: generated.New(db),
	}
}
