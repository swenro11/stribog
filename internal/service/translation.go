package service

import (
	"context"
	"fmt"

	"github.com/swenro11/stribog/internal/entity"
)

// TranslationService -.
type TranslationService struct {
	repo   TranslationRepo
	webAPI TranslationWebAPI
}

// NewTranslationService -.
func NewTranslationService(r TranslationRepo, w TranslationWebAPI) *TranslationService {
	return &TranslationService{
		repo:   r,
		webAPI: w,
	}
}

// History - getting translate history from store.
func (service *TranslationService) History(ctx context.Context) ([]entity.Translation, error) {
	translations, err := service.repo.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("TranslationService - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}

// Translate -.
func (service *TranslationService) Translate(ctx context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := service.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationService - Translate - s.webAPI.Translate: %w", err)
	}

	err = service.repo.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationService - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
