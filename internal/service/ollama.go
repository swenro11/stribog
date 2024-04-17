package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	lingooseOllama "github.com/henomis/lingoose/llm/ollama"
	"github.com/henomis/lingoose/thread"
	"github.com/hupe1980/golc/integration/ollama"
	"github.com/hupe1980/golc/model/llm"
	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
)

const (
	Top13bModel = "wizard-vicuna-uncensored:13b"
)

type OllamaService struct {
	cfg *config.Config
	log *log.Logger
}

func NewOllamaService(cfg *config.Config, l *log.Logger) *OllamaService {
	return &OllamaService{
		cfg: cfg,
		log: l,
	}
}

// based on https://hupe1980.github.io/golc/docs/llms_and_prompts/models/llms/ollama/
// TODO: fix OllamaService - ollamaLlm.Generate: invalid character '{' after top-level value
func (service *OllamaService) Generate(prompt string) (*string, error) {
	client := ollama.New("http://localhost:11434")
	ollamaLlm, errNewOllama := llm.NewOllama(client, func(o *llm.OllamaOptions) {
		o.ModelName = Top13bModel
	})
	if errNewOllama != nil {
		return nil, fmt.Errorf("OllamaService - llm.NewOllama: " + errNewOllama.Error())
	}

	res, errGenerate := ollamaLlm.Generate(context.Background(), prompt)
	if errGenerate != nil {
		return nil, fmt.Errorf("OllamaService - ollamaLlm.Generate: " + errGenerate.Error())
	}

	return &res.Generations[0].Text, nil
}

// based on https://github.com/henomis/lingoose/blob/main/examples/llm/ollama/multimodal/main.go
func (service *OllamaService) GenerateLingoose(prompt string) (*string, error) {
	ollamallm := lingooseOllama.New().WithModel(Top13bModel)

	t := thread.New().AddMessage(
		thread.NewUserMessage().AddContent(
			thread.NewTextContent(prompt),
		),
	)
	/*
		.AddContent(
				thread.NewImageContentFromURL("https://upload.wikimedia.org/wikipedia/commons/thumb/3/34/Anser_anser_1_%28Piotr_Kuczynski%29.jpg/1280px-Anser_anser_1_%28Piotr_Kuczynski%29.jpg"),
			)
	*/

	errGenerate := ollamallm.Generate(context.Background(), t)
	if errGenerate != nil {
		return nil, fmt.Errorf("OllamaService - ollamallm.GenerateLingoose: " + errGenerate.Error())
	}

	strResult := t.String()
	return &strResult, nil
}

func (service *OllamaService) GenerateByPromptWithParam(prompt string, param string) ([]string, error) {
	prompt = fmt.Sprintf(prompt, param)
	prompt += ". Do not include any explanations, only provide a list, every part of list from new row, without numbers."
	result, errGeneratePrompt := service.GenerateLingoose(prompt)
	if errGeneratePrompt != nil {
		return nil, fmt.Errorf("GenerateByPromptWithParam - GenerateLingoose: %s", errGeneratePrompt)
	}

	assistantAnswer := strings.Split(*result, "assistant:\n\tType: text\n\tText:")
	resultStrings := strings.Split(assistantAnswer[1], "\n")
	re := regexp.MustCompile(`\d`)
	for i := range resultStrings {
		fullString := re.ReplaceAllString(resultStrings[i], "")
		fullString = strings.ReplaceAll(fullString, ". ", "")
		resultStrings[i] = strings.TrimSpace(fullString)
	}

	return resultStrings, nil
}

// create content with embedding
// base on https://lingoose.io/reference/embedding/
/*
	https://github.com/henomis/lingoose/blob/main/examples/embeddings/ollama/main.go
	https://github.com/henomis/lingoose/blob/525cbb06fce6b3c2f280374bc0f7dc905eed9f26/examples/embeddings/ollama/main.go#L7
	https://github.com/Burakbgmk/go-tbc-bot/blob/77c0a66e1efe1b2dec8fa146558cedfe8d17a302/internal/ai/query.go#L27
		embeddins, err := ollamaembedder.New().
			WithEndpoint("http://localhost:11434/api").
			WithModel("mistral").
			Embed(
				context.Background(),
				[]string{"What is the NATO purpose?"},
			)
		if err != nil {
			panic(err)
		}
*/
