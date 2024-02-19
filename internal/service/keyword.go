package service

import (
	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ()

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

func (service *KeywordService) CreateKeyword(title string) error {
	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("CreateKeyword - gorm.Open: %s", err)
	}

	db.Create(&entity.Keyword{Title: title, Status: StatusNew})

	return nil
}
