package service

import (
	"context"
	"fmt"

	"github.com/henomis/lingoose/chat"
	"github.com/henomis/lingoose/llm/openai"
	"github.com/henomis/lingoose/prompt"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
)

const (
	/*
		{"object":"list","data":[
			{"id":"gpt-3.5-turbo","object":"model"},
			{"id":"gpt4all-j","object":"model"},
			{"id":"gpt4all-snoozy-13b","object":"model"},
			{"id":"hermes-llama2-13b","object":"model"},
			{"id":"mistral-7b-openorca.Q4_0.gguf","object":"model"}
		]}
	*/
	_LocalAIgpt3dot5turbo = "gpt-3.5-turbo"
)

type LocalAIService struct {
	cfg *config.Config
	log *log.Logger
}

func NewLocalAIService(cfg *config.Config, l *log.Logger) *LocalAIService {
	return &LocalAIService{
		cfg: cfg,
		log: l,
	}
}

// based on https://github.com/henomis/lingoose/blob/main/examples/llm/openai/localai/main.go
func (service *LocalAIService) TextGenerationGpt3dot5turbo(promt string) (string, error) {

	customConfig := goopenai.DefaultConfig("")
	customConfig.BaseURL = service.cfg.AI.LocalAIURL
	customClient := goopenai.NewClientWithConfig(customConfig)

	//openaiModel := new Model
	chat := chat.New(
		chat.PromptMessage{
			Type:   chat.MessageTypeUser,
			Prompt: prompt.New(promt),
		},
	)

	llm := openai.NewChat().WithClient(customClient).WithModel(openai.GPT3Dot5Turbo)

	result, errCompletion := llm.Chat(context.Background(), chat)
	if errCompletion != nil {
		return result, fmt.Errorf("LocalAIService.TextGenerationGpt3dot5turbo - llm.Chat: " + errCompletion.Error())
	}

	return result, nil
}
