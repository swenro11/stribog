package service

import (
	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	StatusNew              string = "New"
	StatusApprove          string = "Approve"
	StatusGenerating       string = "Generating"      //StartGegenerateImages && SaveImages
	StatusReadyWithImages  string = "ReadyWithImages" //SeoOptimization
	StatusCheckUnique      string = "CheckUnique"
	StatusCheckNotAI       string = "CheckNotAI"
	StatusDeployToTestHugo string = "DeployToTestHugo" //MdFormatWithImages
	StatusRewrite          string = "Rewrite"
	StatusDeployToProdHugo string = "DeployToProdHugo"
	//StatusAddLinks            = "AddLinks"
	//StatusRegenerateImages    = "RegenerateImages"
	//StatusGenerateTags        = "GenerateTags"
)

type ArticleService struct {
	cfg *config.Config
	log *log.Logger
}

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

	db.Create(&entity.Article{Title: title, Status: StatusNew})

	return nil
}

func (service *ArticleService) CreateArticleWithImages(keyword string) error {
	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	db.Create(&entity.Article{Title: keyword, Status: StatusNew})

	//create []Image with generaten promts & StatusNew

	return nil
}

/*
ArticleService. Hugo
ArticleService. Images
ArticleService. CheckNotAI - https://www.zerogpt.com/ and https://gptzero.me/ (https://github.com/BurhanUlTayyab/GPTZero)
ArticleService. CheckUnique - based on https://plagiarismcheck.org/for-developers/#single
ArticleService. SeoOptimization
*/
