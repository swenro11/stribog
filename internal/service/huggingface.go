package service

import (
	"context"
	"fmt"
	"os"

	"github.com/henomis/lingoose/llm/huggingface"
	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
)

const (
	_HuggingfaceTokenParam = "HUGGING_FACE_HUB_TOKEN"
)

// HuggingfaceService -.
type HuggingfaceService struct {
	cfg *config.Config
	log *log.Logger
}

// NewHuggingfaceService -.
func NewHuggingfaceService(cfg *config.Config, l *log.Logger) *HuggingfaceService {
	return &HuggingfaceService{
		cfg: cfg,
		log: l,
	}
}

func (service *HuggingfaceService) checkGetenv(enableLog bool) string {
	huggingfaceToken := os.Getenv(_HuggingfaceTokenParam)
	if enableLog {
		service.log.Info("HuggingfaceService - checkGetenv -  os.Getenv = " + huggingfaceToken)
	}
	if len(huggingfaceToken) > 0 {
		return huggingfaceToken
	}

	os.Setenv(_HuggingfaceTokenParam, service.cfg.PARAM.HuggingfaceToken)
	huggingfaceToken = os.Getenv(_HuggingfaceTokenParam)
	if enableLog {
		service.log.Info("HuggingfaceService - checkGetenv -  second os.Getenv = " + huggingfaceToken)
	}

	return huggingfaceToken
}

func (service *HuggingfaceService) TextGeneration(enableLog bool, ctx context.Context, model string, promt string) (string, error) {

	llm := huggingface.New(model, 0.1, false).WithMode(huggingface.ModeTextGeneration)

	result, errCompletion := llm.Completion(context.Background(), promt)
	if errCompletion != nil {
		return result, fmt.Errorf("HuggingfaceService - llm.Completion: " + errCompletion.Error())
	}

	return result, nil
	/*
		_, err = llm.BatchCompletion(
			context.Background(),
			[]string{
				"Write a joke about geese.",
				"What is the NATO purpose?",
			},
		)
		if err != nil {
			panic(err)
		}
	*/
}
