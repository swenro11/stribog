package repo

import (
	"github.com/swenro11/stribog/pkg/postgres"
)

// PoolRepo -.
type PoolRepo struct {
	*postgres.Postgres
}

// PoolRepo -.
func NewPoolRepo(pg *postgres.Postgres) *PoolRepo {
	return &PoolRepo{pg}
}
