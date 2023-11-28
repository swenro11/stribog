// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"

	entity "github.com/swenro11/stribog/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=service_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// TokenInfoRepo -.
	TokenInfoRepo interface {
		GetAllTokens(context.Context) ([]entity.TokenInfo, error)
		GetBlockchainTokens(context.Context, string) ([]entity.TokenInfo, error)
		GetByIdAndType(context.Context, string, string) ([]entity.TokenInfo, error)
		StoreTokenInfo(context.Context, entity.TokenInfo) error
		UpdateTokenInfo(context.Context, entity.TokenInfo, entity.TokenInfo) error
		DeleteAllByBlockchain(context.Context, string) error
	}

	// CgCoinRepo -.
	CgCoinRepo interface {
		GetAllCoins(context.Context) ([]entity.CgCoin, error)
		StoreCgCoin(context.Context, entity.CgCoin) error
		DeleteAll(context.Context) error
	}

	// Tasks -.
	Tasks interface {
		CheckRabbitTask() string
	}
)
