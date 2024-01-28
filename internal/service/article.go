package service

import (
	//"context"
	//"fmt"

	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	_StatusNew                   = "New"
	_StatusSeoOptimization       = "SeoOptimization"
	_StatusCheckUnique           = "CheckUnique"
	_StatusCheckNotAI            = "CheckNotAI"
	_StatusStartGegenerateImages = "StartGegenerateImages"
	_StatusSaveImages            = "SaveImages"
	_StatusDeployToTestHugo      = "DeployToTestHugo" //MdFormatWithImages
	_StatusRewrite               = "Rewrite"
	_StatusDeployToProdHugo      = "DeployToProdHugo"
	//_StatusAddLinks            = "AddLinks"
	//_StatusRegenerateImages    = "RegenerateImages"
	//_StatusGenerateTags        = "GenerateTags"
)

// ArticleService -.
type ArticleService struct {
	cfg *config.Config
	log *log.Logger
}

// NewArticleService -.
func NewArticleService(cfg *config.Config, l *log.Logger) *ArticleService {
	return &ArticleService{
		cfg: cfg,
		log: l,
	}
}

func (service *ArticleService) CreateArticle(title string) error {
	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	db.Create(&entity.Article{Title: title, Status: _StatusNew})

	return nil
}
