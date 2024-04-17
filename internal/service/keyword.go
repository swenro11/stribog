package service

import (
	"fmt"
	"strings"

	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	//SystemPrompt = "" ? i need that? from https://community.openai.com/t/getting-response-data-as-a-fixed-consistent-json-response/28471/22?page=2
	KeywordsPrompt          = "Generate a list of at least 10 keywords related to %s"
	ClusterKeywordsPrompt   = "Generate a cluster of keywords around the primary keyword '%s'"
	KeywordVariationsPrompt = "Create keyword variations for '%s' with high search volume"
	BukvarixSource          = "Bukvarix"
	CohereSource            = "Cohere"
	OllamaSource            = "Ollama"
)

type KeywordService struct {
	cfg *config.Config
	log *log.Logger
}

func NewKeywordService(cfg *config.Config, l *log.Logger) *KeywordService {
	return &KeywordService{
		cfg: cfg,
		log: l,
	}
}

func (service *KeywordService) BukvarixSaveKeywords(topic entity.Topic) error {
	bukvarixService := NewBukvarixService(
		service.cfg,
		service.log,
	)

	resultKeywords, errKeywords := bukvarixService.Keywords(topic.Title)
	if errKeywords != nil {
		return fmt.Errorf("KeywordService.BukvarixSaveKeywords - bukvarixService.Keywords: %s", errKeywords)
	}

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("KeywordService.CreateKeywords - gorm.Open: %s", err)
	}

	for _, element := range resultKeywords {
		if element != "" {
			db.Create(&entity.Keyword{TopicID: topic.ID, Status: StatusNew, Source: BukvarixSource, Title: element})
		}
	}

	return nil
}

// TODO: test
func (service *KeywordService) CohereSaveLongtailKeywords(topic entity.Topic) error {
	cohereService := NewCohereService(
		service.cfg,
		service.log,
	)

	prompt := "Generate a list of long-tail keywords related to '" + topic.Title + "'"
	prompt += ".Do not include any explanations, only provide a list with keywords, in cvs format with ; separator."
	result, errGeneratePrompt := cohereService.GeneratePrompt(prompt)
	if errGeneratePrompt != nil {
		return fmt.Errorf("KeywordService.CohereSaveLongtailKeywords - cohereService.GeneratePrompt: %s", errGeneratePrompt)
	}

	service.log.Info("cohereService.GeneratePrompt = " + *result)

	/*
		service.log.Info("cohereService.GeneratePrompt = " + *result)
		Certainty; Belief; Consciousness; Existentialism; Freedom; Knowledge; Mind; Morality; Nature of God; Science and Philosophy; Suffering
	*/

	resultKeywords := strings.Split(*result, ";")

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("KeywordService.CohereSaveLongtailKeywords - gorm.Open: %s", err)
	}

	for _, element := range resultKeywords {
		if element != "" {
			title := strings.Trim(element, " ")
			db.Create(&entity.Keyword{TopicID: topic.ID, Status: StatusNew, Source: CohereSource, Title: title})
		}
	}

	return nil
}

// TODO: Use for new topic creatinon
// TODO: Compare with Example prompt: Generate a cluster of keywords around the primary keyword "blockchain technology."
func (service *KeywordService) CohereSaveKeywords(topic entity.Topic) error {
	cohereService := NewCohereService(
		service.cfg,
		service.log,
	)

	prompt := "Generate a list of at least 10 keywords related to " + topic.Title
	prompt += ".Do not include any explanations, only provide a list with keywords, in cvs format with ; separator."
	result, errGeneratePrompt := cohereService.GeneratePrompt(prompt)
	if errGeneratePrompt != nil {
		return fmt.Errorf("KeywordService.CohereSaveKeywords - cohereService.GeneratePrompt: %s", errGeneratePrompt)
	}

	/*
		service.log.Info("cohereService.GeneratePrompt = " + *result)
		Certainty; Belief; Consciousness; Existentialism; Freedom; Knowledge; Mind; Morality; Nature of God; Science and Philosophy; Suffering
	*/

	resultKeywords := strings.Split(*result, ";")

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("KeywordService.CohereSaveKeywords - gorm.Open: %s", err)
	}

	for _, element := range resultKeywords {
		if element != "" {
			title := strings.Trim(element, " ")
			db.Create(&entity.Keyword{TopicID: topic.ID, Status: StatusNew, Source: CohereSource, Title: title})
		}
	}

	return nil
}

func (service *KeywordService) OllamaSaveKeywords(prompt string, topic entity.Topic) error {
	ollamaService := NewOllamaService(
		service.cfg,
		service.log,
	)

	resultKeywords, errGenerateByPromptWithParam := ollamaService.GenerateByPromptWithParam(prompt, topic.Title)
	if errGenerateByPromptWithParam != nil {
		return fmt.Errorf("KeywordService.OllamaSaveKeywords - ollamaService.GenerateByPromptWithParam: %s", errGenerateByPromptWithParam)
	}

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("KeywordService.OllamaSaveKeywords - gorm.Open: %s", err)
	}

	for _, element := range resultKeywords {
		if element != "" {
			db.Create(&entity.Keyword{TopicID: topic.ID, Status: StatusNew, Source: OllamaSource, Title: element})
		}
	}

	return nil
}
