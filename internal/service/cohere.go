package service

import (
	"context"
	"fmt"

	"github.com/henomis/lingoose/llm/cohere"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
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

func (service *CohereService) Completion(prompt string) (*string, error) {

	llm := cohere.NewCompletion().WithAPIKey(service.cfg.AI.CohereToken).WithMaxTokens(100).WithTemperature(0.1).WithVerbose(true)

	result, errCompletion := llm.Completion(context.Background(), prompt)
	if errCompletion != nil {
		return nil, fmt.Errorf("CohereService - llm.Completion: " + errCompletion.Error())
	}

	return &result, nil
}

// based on https://hupe1980.github.io/golc/docs/llms_and_prompts/models/chatmodels/cohere/
func (service *CohereService) GeneratePrompt(inputPrompt string) (*string, error) {
	cohere, errNewCohere := chatmodel.NewCohere(service.cfg.AI.CohereToken, func(o *chatmodel.CohereOptions) {
		o.Temperature = 0.7 // optional
	})
	if errNewCohere != nil {
		return nil, fmt.Errorf("CohereService - chatmodel.NewCohere: " + errNewCohere.Error())
	}

	res, errGeneratePrompt := model.GeneratePrompt(context.Background(), cohere, prompt.StringPromptValue(inputPrompt))
	if errGeneratePrompt != nil {
		return nil, fmt.Errorf("CohereService - model.GeneratePrompt: " + errGeneratePrompt.Error())
	}

	return &res.Generations[0].Text, nil
}
