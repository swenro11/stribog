package service

import (
	"context"
	"fmt"

	"github.com/henomis/lingoose/llm/cohere"
	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
)

type CohereService struct {
	cfg *config.Config
	log *log.Logger
}

func NewCohereService(cfg *config.Config, l *log.Logger) *CohereService {
	return &CohereService{
		cfg: cfg,
		log: l,
	}
}

func (service *CohereService) TextGeneration(promt string) (string, error) {

	llm := cohere.NewCompletion().WithAPIKey(service.cfg.PARAM.CohereToken).WithMaxTokens(100).WithTemperature(0.1).WithVerbose(true)

	result, errCompletion := llm.Completion(context.Background(), promt)
	if errCompletion != nil {
		return result, fmt.Errorf("CohereService - llm.Completion: " + errCompletion.Error())
	}

	return result, nil
}
