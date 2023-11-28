package repo

import (
	"context"
	"fmt"

	"github.com/swenro11/stribog/internal/entity"
	"github.com/swenro11/stribog/pkg/postgres"
)

// CgCoinRepo -.
type CgCoinRepo struct {
	*postgres.Postgres
}

// NewCgCoinRepo -.
func NewCgCoinRepo(pg *postgres.Postgres) *CgCoinRepo {
	return &CgCoinRepo{pg}
}

// GetAllCoins -.
func (r *CgCoinRepo) GetAllCoins(ctx context.Context) ([]entity.CgCoin, error) {
	sql, _, err := r.Builder.
		Select("id",
			"symbol",
			"name").
		From("cgcoin").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CgCoinRepo - GetAllCoins - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("CgCoinRepo - GetAllCoins - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.CgCoin, 0, _defaultEntityCap)

	for rows.Next() {
		m := entity.CgCoin{}

		err = rows.Scan(
			&m.ID,
			&m.Symbol,
			&m.Name)
		if err != nil {
			return nil, fmt.Errorf("CgCoinRepo - GetAllCoins - rows.Scan: %w", err)
		}

		entities = append(entities, m)
	}

	return entities, nil
}

// StoreCgCoin -.
func (r *CgCoinRepo) StoreCgCoin(ctx context.Context, m entity.CgCoin) error {
	sql, args, err := r.Builder.
		Insert("cgcoin").
		Columns(
			"id",
			"symbol",
			"name").
		Values(
			m.ID,
			m.Symbol,
			m.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("CgCoinRepo - StoreCgCoin - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("CgCoinRepo - StoreCgCoin - r.Pool.Exec: %w", err)
	}

	return nil
}

// DeleteAll -.
func (r *CgCoinRepo) DeleteAll(ctx context.Context) error {
	sql, args, err := r.Builder.
		Delete("cgcoin").
		ToSql()
	if err != nil {
		return fmt.Errorf("CgCoinRepo - DeleteAll - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("CgCoinRepo - DeleteAll - r.Pool.Exec: %w", err)
	}

	return nil
}
