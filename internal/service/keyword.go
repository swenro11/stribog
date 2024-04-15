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
	GenerateKeywordsPrompt = "Generate a list of at least 10 keywords related to " //[topic]
	BukvarixSource         = "Bukvarix"
	CohereSource           = "Cohere"
	OllamaSource           = "Ollama"
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

func (service *KeywordService) CohereSaveKeywords(topic entity.Topic) error {
	cohereService := NewCohereService(
		service.cfg,
		service.log,
	)

	prompt := "Generate a list of at least 10 keywords related to " + topic.Title
	prompt += ".Do not include any explanations, only provide a list with keywords, in cvs format with ; separator."
	result, errGeneratePrompt := cohereService.GeneratePrompt(prompt)
	if errGeneratePrompt != nil {
		return fmt.Errorf("KeywordService.CreateKeywords - cohereService.GeneratePrompt: %s", errGeneratePrompt)
	}

	service.log.Info("KeywordService.GeneratePrompt = " + *result)

	/* result
	Certainty; Belief; Consciousness; Existentialism; Freedom; Knowledge; Mind; Morality; Nature of God; Science and Philosophy; Suffering
	*/

	resultKeywords := strings.Split(*result, ";")

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("KeywordService.CreateKeywords - gorm.Open: %s", err)
	}

	for _, element := range resultKeywords {
		if element != "" {
			title := strings.Trim(element, " ")
			db.Create(&entity.Keyword{TopicID: topic.ID, Status: StatusNew, Source: CohereSource, Title: title})
		}
	}

	return nil
}

/*
prompt := "Generate a list of at least 10 keywords related to " + topic.Title
prompt += ".Do not include any explanations, only provide a RFC8259 compliant JSON response following this format without deviation."
prompt += "['keyword one', 'keyword two', 'etc.']"
result, errGeneratePrompt := cohereService.GeneratePrompt(prompt)
if errGeneratePrompt != nil {
	return fmt.Errorf("KeywordService.CreateKeywords - cohereService.GeneratePrompt: %s", errGeneratePrompt)
}

service.log.Info("KeywordService.GeneratePrompt = " + *result)
result
Here is a list of 10 keywords related to Life Philosophy in RFC8259 compliant JSON response format:\n\n```json\n[\n \"life\",\n \"existence\",\n \"meaning\",\n \"purpose\",\n \"values\",\n \"truth\",\n \"awareness\",\n \"authenticity\",\n \"survival\",\n \"self-discovery\"\n]\n``` \n\nThese keywords were chosen after crawling extensively through texts, discourses, and thoughts from various different philosophers across time.
*/
