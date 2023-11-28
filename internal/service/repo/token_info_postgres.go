package repo

import (
	"context"
	"fmt"

	"github.com/swenro11/stribog/internal/entity"
	"github.com/swenro11/stribog/pkg/postgres"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/ethclient"
)

// TokenInfoRepo -.
type TokenInfoRepo struct {
	*postgres.Postgres
}

/*
type Erc20 interface {
	LogoUrl() string
	Code() string
	Name() string
	Address() string
	BalanceOf(client *ethclient.Client, account common.Address) (uint64, error)
	FormattedBalance(client *ethclient.Client, account common.Address) (float64, error)
}
*/

// TokenInfoRepo -.
func NewTokenInfoRepo(pg *postgres.Postgres) *TokenInfoRepo {
	return &TokenInfoRepo{pg}
}

// GetAllTokens -.
func (r *TokenInfoRepo) GetAllTokens(ctx context.Context) ([]entity.TokenInfo, error) {
	sql, _, err := r.Builder.
		Select("id",
			"name",
			"website",
			"description",
			"explorer",
			"type",
			"symbol",
			"status",
			"decimals").
		From("token_info").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TokenInfoRepo - GetAllTokens - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TokenInfoRepo - GetAllTokens - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.TokenInfo, 0, _defaultEntityCap)

	for rows.Next() {
		m := entity.TokenInfo{}

		err = rows.Scan(
			&m.Id,
			&m.Name,
			&m.Website,
			&m.Description,
			&m.Explorer,
			&m.Type,
			&m.Symbol,
			&m.Status,
			&m.Decimals)
		if err != nil {
			return nil, fmt.Errorf("TokenInfoRepo - GetAllTokens - rows.Scan: %w", err)
		}

		entities = append(entities, m)
	}

	return entities, nil
}

// GetBlockchainTokens -.
func (r *TokenInfoRepo) GetBlockchainTokens(ctx context.Context, blockchain string) ([]entity.TokenInfo, error) {
	sql, _, err := r.Builder.
		Select("id",
			"name",
			"website",
			"description",
			"explorer",
			"type",
			"symbol",
			"status",
			"decimals").
		From("token_info").
		Where("type = '" + blockchain + "'").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TokenInfoRepo - GetBlockchainTokens - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TokenInfoRepo - GetBlockchainTokens - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.TokenInfo, 0, _defaultEntityCap)

	for rows.Next() {
		m := entity.TokenInfo{}

		err = rows.Scan(
			&m.Id,
			&m.Name,
			&m.Website,
			&m.Description,
			&m.Explorer,
			&m.Type,
			&m.Symbol,
			&m.Status,
			&m.Decimals)
		if err != nil {
			return nil, fmt.Errorf("TokenInfoRepo - GetBlockchainTokens - rows.Scan: %w", err)
		}

		entities = append(entities, m)
	}

	return entities, nil
}

// GetBlockchainTokens -.
func (r *TokenInfoRepo) GetByIdAndType(ctx context.Context, blockchain string, id string) ([]entity.TokenInfo, error) {
	sql, _, err := r.Builder.
		Select("id",
			"name",
			"website",
			"description",
			"explorer",
			"type",
			"symbol",
			"status",
			"decimals").
		From("token_info").
		Where("type = '" + blockchain + "' and id = '" + id + "'").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TokenInfoRepo - GetByIdAndType - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TokenInfoRepo - GetByIdAndType - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.TokenInfo, 0, _defaultEntityCap)

	for rows.Next() {
		m := entity.TokenInfo{}

		err = rows.Scan(
			&m.Id,
			&m.Name,
			&m.Website,
			&m.Description,
			&m.Explorer,
			&m.Type,
			&m.Symbol,
			&m.Status,
			&m.Decimals)
		if err != nil {
			return nil, fmt.Errorf("TokenInfoRepo - GetByIdAndType - rows.Scan: %w", err)
		}

		entities = append(entities, m)
	}

	return entities, nil
}

// StoreTokenInfo -.
func (r *TokenInfoRepo) StoreTokenInfo(ctx context.Context, t entity.TokenInfo) error {
	sql, args, err := r.Builder.
		Insert("token_info").
		Columns("id",
			"name",
			"website",
			"description",
			"explorer",
			"type",
			"symbol",
			"status",
			"decimals").
		Values(
			t.Id,
			t.Name,
			t.Website,
			t.Description,
			t.Explorer,
			t.Type,
			t.Symbol,
			t.Status,
			t.Decimals).
		ToSql()
	if err != nil {
		return fmt.Errorf("TokenInfoRepo - StoreTokenInfo - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TokenInfoRepo - StoreTokenInfo - r.Pool.Exec: %w", err)
	}

	return nil
}

// UpdateTokenInfo -.
func (r *TokenInfoRepo) UpdateTokenInfo(ctx context.Context, existToken entity.TokenInfo, tokenInfo entity.TokenInfo) error {
	sql, args, err := r.Builder.
		Update("token_info").
		Set("name", tokenInfo.Name).
		Set("website", tokenInfo.Website).
		Set("description", tokenInfo.Description).
		Set("explorer", tokenInfo.Explorer).
		Set("symbol", tokenInfo.Symbol).
		Set("status", tokenInfo.Status).
		Set("decimals", tokenInfo.Decimals).
		Where("id = '" + tokenInfo.Id + "' and type = '" + tokenInfo.Type + "'").
		ToSql()
	if err != nil {
		return fmt.Errorf("TokenInfoRepo - UpdateTokenInfo - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TokenInfoRepo - UpdateTokenInfo - r.Pool.Exec: %w", err)
	}

	return nil
}

// DeleteAllByBlockchain -.
func (r *TokenInfoRepo) DeleteAllByBlockchain(ctx context.Context, blockchain string) error {
	sql, args, err := r.Builder.
		Delete("token_info").
		Where("type = '" + blockchain + "'").
		ToSql()
	if err != nil {
		return fmt.Errorf("TokenInfoRepo - DeleteAllByBlockchain - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TokenInfoRepo - DeleteAllByBlockchain - r.Pool.Exec: %w", err)
	}

	return nil
}

/*
func (r *TokenInfoRepo) NewTokenFromERC20(t Erc20) *entity.TokenInfo {
	return &entity.TokenInfo{
		Address:     t.Address(),
		Name:        t.Name(),
		Symbol:      t.Code(),
		LogoUrl:     t.LogoUrl(),
		Decimals:    0,
		TotalSupply: 0,
		LoanAmount:  0,
	}
}
*/
