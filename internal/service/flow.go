package service

import (
	"fmt"

	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
// BaseURL              = "https://api-key.fusionbrain.ai/key/api/v1/"
)

type FlowService struct {
	cfg *config.Config
	log *log.Logger
}

// example
type Flow struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

func NewFlowService(cfg *config.Config, l *log.Logger) *FlowService {
	return &FlowService{
		cfg: cfg,
		log: l,
	}
}

/*
- DeployToProdHugo
*/
func (service *FlowService) DeployToProdHugo() (*string, error) {
	// Mock

	return nil, nil
}

/*
- Save MdFormatWithImages to filesystem
- Create new file in current project folder ($slugFilename, $category) : string (path/to/file) | error
hugo new posts/my-first-post.md
*/
func (service *FlowService) DeployToTestHugo() (*string, error) {
	// Mock

	return nil, nil
}

/*
- Start SeoOptimization
- CheckUnique
- CheckNotAI
*/
func (service *FlowService) Rewrite() (*string, error) {
	// Mock

	return nil, nil
}

/*
- Start SeoOptimization
- CheckUnique -> Rewrite Or DeployToTestHugo
- finally CheckNotAI -> Rewrite Or DeployToTestHugo
*/
func (service *FlowService) ReadyWithImages() (*string, error) {
	// Mock

	return nil, nil
}

/*
-
*/
func (service *FlowService) ApprovedByAI() (*string, error) {
	// Mock

	return nil, nil
}

/*
- StartGegenerateImages && SaveImages
*/
func (service *FlowService) Generating() (*string, error) {
	// Mock
	// Mb take a part of Approved for fussiabbrain here

	return nil, nil
}

/*
- Create articles for Approved Keywords
- Start creating images
- Get & save images from fusionbrain
*/
func (service *FlowService) Approved() (*string, error) {
	WriterService := NewWriterService(
		service.cfg,
		service.log,
	)

	fusionbrainService := NewFusionbrainService(
		service.cfg,
		service.log,
	)

	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	var keywords []entity.Keyword
	db.Where(entity.Keyword{Status: StatusApproved}).Find(&keywords)
	for _, keyword := range keywords {
		errCreateArticle := WriterService.CreateArticleWithImages(keyword)
		if errCreateArticle != nil {
			return nil, fmt.Errorf("FlowService.Approved - WriterService.CreateArticleWithImages: %s", errCreateArticle)
		}
	}
	//fusionbrain online
	_, errGetModels := fusionbrainService.GetModels()
	if errGetModels == nil {
		var images []entity.Image
		db.Where(entity.Image{Status: StatusNew}).Find(&images)
		for _, image := range images {
			_, errCreateTask := fusionbrainService.CreateTaskForImage(image, 1024, 1024, "", "", false)
			if errCreateTask != nil {
				return nil, fmt.Errorf("FlowService.Approved - fusionbrainService.CreateTaskForImage: %s", errCreateTask)
			}
		}

		var tasks []entity.Task
		db.Where(entity.Keyword{Status: TaskStatusInitial}).Find(&tasks)
		for _, task := range tasks {
			service.log.Info("task.Uuid: %s", task.Uuid)
			getImagesResult, errGetImages := fusionbrainService.GetImages(&task, false)
			if errGetImages != nil {
				return nil, fmt.Errorf("FlowService.Approved - fusionbrainService.CreateTaskForImage: %s", errGetImages)
			}
			service.log.Info(getImagesResult.Uuid)
		}
	}

	return nil, nil
}

/*
- Create New keywords by LLM & other
- Create topics from sources
- Use current article table as a source for new topics
*/
func (service *FlowService) New() (*string, error) {
	// Mock Create New keywords by LLM & other
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
		errSaveKeyword := keywordService.BukvarixSaveKeywords(topic)
		if errSaveKeyword != nil {
			service.log.Fatal(errSaveKeyword.Error())
		}

		errCreateKeyword := keywordService.CohereCreateKeywords(topic)
		if errCreateKeyword != nil {
			service.log.Fatal(errCreateKeyword.Error())
		}
	}

	// TODO: Move to OllamaService
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

	return nil, nil
}
