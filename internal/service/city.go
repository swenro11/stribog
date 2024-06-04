package service

import (
	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CityService struct {
	cfg *config.Config
	log *log.Logger
}

func NewCityService(cfg *config.Config, l *log.Logger) *CityService {
	return &CityService{
		cfg: cfg,
		log: l,
	}
}

func (service *CityService) New() (*string, error) {
	keywordService := NewKeywordService(
		service.cfg,
		service.log,
	)

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	var topics []entity.Topic
	db.Where(entity.Image{Status: StatusApproved}).Find(&topics)
	for _, topic := range topics {
		errCohereSaveKeywords := keywordService.CohereSaveKeywords(topic)
		if errCohereSaveKeywords != nil {
			service.log.Fatal(errCohereSaveKeywords.Error())
		}

		errSaveKeyword := keywordService.OllamaSaveKeywords(KeywordsPrompt, topic)
		if errSaveKeyword != nil {
			service.log.Fatal(errSaveKeyword.Error())
		}
	}

	return nil, nil
}
