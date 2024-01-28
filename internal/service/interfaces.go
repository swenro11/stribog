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

	// Tasks -.
	Tasks interface {
		CheckRabbitTask() string
	}

	// PoolRepo -.
	PoolRepo interface {
	}
)
