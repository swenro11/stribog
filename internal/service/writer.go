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

type WriterService struct {
	cfg *config.Config
	log *log.Logger
}

func NewWriterService(cfg *config.Config, l *log.Logger) *WriterService {
	return &WriterService{
		cfg: cfg,
		log: l,
	}
}

func (service *WriterService) CreateArticle(title string) error {
	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	db.Create(&entity.Article{Title: title, Status: StatusNew})

	return nil
}

func (service *WriterService) CreateArticleWithImages(keyword string) error {
	// TODO: generate Seo title from keyword
	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	db.Create(&entity.Article{Title: keyword, Status: StatusNew})

	//create []Image with generation promts & StatusNew

	return nil
}

/*
WriterService. Hugo
WriterService. Images
WriterService. CheckNotAI - https://www.zerogpt.com/ and https://gptzero.me/ (https://github.com/BurhanUlTayyab/GPTZero)
WriterService. CheckUnique - based on https://plagiarismcheck.org/for-developers/#single
WriterService. SeoOptimization
*/
