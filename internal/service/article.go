package service

import (
	//"context"
	//"fmt"

	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
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
