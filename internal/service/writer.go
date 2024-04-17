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
	StatusApproved         string = "Approved"
	StatusApprovedByAI     string = "ApprovedByAI"
	StatusGenerating       string = "Generating"
	StatusReadyWithImages  string = "ReadyWithImages"
	StatusCheckUnique      string = "CheckUnique"
	StatusCheckNotAI       string = "CheckNotAI" // https://www.zerogpt.com/ and https://gptzero.me/ (https://github.com/BurhanUlTayyab/GPTZero)
	StatusDeployToTestHugo string = "DeployToTestHugo"
	StatusRewrite          string = "Rewrite"
	StatusDeployToProdHugo string = "DeployToProdHugo"
	//StatusAddLinks            = "AddLinks"
	//StatusRegenerateImages    = "RegenerateImages"
	//StatusGenerateTags        = "GenerateTags"
	//CheckUnique - based on https://plagiarismcheck.org/for-developers/#single
	MetaDescriptionPrompt       = "Write a meta description of 150-160 characters for a blog post on '%s'"
	IntroductionParagraphPrompt = "Generate an introduction paragraph for an article about '%s'"
	SEOoptimizedTitlePrompt     = "Generate a compelling and SEO-optimized blog title for an article about '%s'"
	ImageAltTagPrompt           = "Curate informative image alt tags for a photo gallery of '%s'"
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

	/*
		huggingfaceService := NewHuggingfaceService(
			cfg,
			service.log,
		)

		result, errHermesTextGen := huggingfaceService.HermesTextGen("What is the NATO purpose?")
		if errHermesTextGen != nil {
			service.log.Fatal(errHermesTextGen)
		}

		result, errTextGeneration := huggingfaceService.FusionNetTextGenHupe1980("What is the NATO purpose?")
		if errTextGeneration != nil {
			service.log.Fatal(errTextGeneration)
		}
	*/

	/*
		localaiService := NewLocalAIService(
			cfg,
			service.log,
		)

		result, errTextGeneration := localaiService.TextGenerationGpt3dot5turbo("What is the NATO purpose?")
		if errTextGeneration != nil {
			service.log.Fatal(errTextGeneration)
		}

		service.log.Info("Result - ", result)
	*/

	/*
		ollamaService := NewOllamaService(
			cfg,
			service.log,
		)

		result, errGenerate := ollamaService.GenerateLingoose("Create seo optimized article about stoicism")
		if errGenerate != nil {
			service.log.Fatal(errGenerate)
		}
		service.log.Info(*result)
		// Thread:\nuser:\n\tType: text\n\tText: Create seo optimized article about stoicism\nassistant:\n\tType: text\n\tText: Stoicism is a philosophy that originated in ancient Greece and has been practiced for thousands of years. It emphasizes the development of self-control, fortitude, and resilience in the face of adversity. The principles of Stoicism can be applied to modern life to help individuals navigate challenges and find inner peace.\n\nOne of the core tenets of Stoicism is the idea that we should focus on things that are within our control, rather than worrying about external factors that are beyond our power to change. This means accepting what has happened in the past and not dwelling on it, as well as being prepared for whatever may come in the future.\n\nAnother important aspect of Stoicism is the practice of mindfulness or being present in the moment. By focusing on the present rather than ruminating about the past or worrying about the future, individuals can cultivate a sense of calm and clarity that can help them navigate challenges more effectively.\n\nIn addition to these principles, Stoicism also emphasizes the importance of self-reflection and personal development. By regularly examining our thoughts and actions, we can identify areas for improvement and work towards becoming better versions of ourselves.\n\nOverall, Stoicism offers a powerful framework for living a meaningful and fulfilling life. By embracing its principles and practicing mindfulness and self-reflection, individuals can cultivate resilience in the face of adversity and find inner peace even in the midst of chaos.\n
	*/

	/*
		//Use VPN
		cohereService := NewCohereService(
			cfg,
			service.log,
		)

		result, errGeneratePrompt := cohereService.GeneratePrompt("Create seo optimized article about stoicism")
		if errGeneratePrompt != nil {
			service.log.Fatal(errGeneratePrompt)
		}
		service.log.Info(*result)

		// Title: Embracing Resilience and Virtue: Unlocking the Secrets of Stoic Philosophy\n\nIntroduction:\nAre you feeling overwhelmed by the uncertainties of life? Discover the timeless wisdom of Stoic philosophy and unlock the secrets to embracing adversity, cultivating resilience, and living a life of virtue. Stoicism holds the key to navigating the challenges of the modern world with clarity, resilience, and self-improvement. Get ready to embark on a transformative journey through the pillars of Stoic philosophy.\n\nPhilosophers of the Stoic Tradition:\nFrom the early days of Zeno of Citium to the influence of Roman Stoics like Seneca and Marcus Aurelius, this philosophy has stood the test of time. Let's delve into their teachings and explore how they can impact our modern lives.\n\n1. Embrace Adversity:\nStoics reject the belief that external events or emotions solely determine one's happiness. Instead, they emphasize the power of our perception and response to adversity. By adopting a rational and pragmatic perspective, individuals can develop resilience in the face of difficulties. \n\n2. Focus on Self-Improvement:\nStoicism emphasizes self-improvement and accountability. Stoics believe that individuals have control over their thoughts, actions, and reactions in any given situation. By embracing self-discipline and practicing introspection, individuals can cultivate personal growth and become better versions of themselves. \n\n3. Rationality as a Guide:\nStoics highly value rationality as a guiding force in life. They encourage individuals to assess things and events in a logical, unbiased manner, rather than letting emotions dictate their actions. By fostering rational thinking, individuals can make better decisions, overcome impulsive behavior, and find peace in clarity. \n\n4. Virtue and Ethics:\nStoic philosophy intertwines with virtue ethics, emphasizing the development of positive character traits. Stoics seek to live in harmony with universal principles and consistently act with courage, wisdom, justice, and moderation. By doing so, they strengthen their moral character and contribute to the common good. \n\n5. Finding Peace in Simplicity:\nStoicism promotes finding tranquility in simplicity. By detaching oneself from material possessions and transient emotions, individuals can discover inner peace and resilience. Simplicity empowers individuals to find contentment and embrace the inherent challenges of the human condition. \n\nApplication in the Modern World:\nLet's explore how we can apply these ancient teachings in our modern lives. Stoicism offers valuable guidance for navigating relationships, career challenges, and the constant pressures of the digital world. It can help us maintain a sense of perspective, make informed decisions, and cultivate emotional well-being, despite the chaos of the world. \n\nConclusion:\nStoic philosophy invites us to challenge our perceptions, cultivate resilience, and embrace a life of virtue and purpose. By applying these timeless principles, we can strengthen our resolve, make informed decisions, and find tranquility in our increasingly complex world. Are you ready to unlock the secrets of stoicism and embark on a journey of self-discovery and resilience?
	*/

	db.Create(&entity.Article{Title: title, Status: StatusNew})

	return nil
}

func (service *WriterService) CreateArticleWithImages(keyword entity.Keyword) error {
	// TODO: generate Seo title from keyword
	db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
	if err != nil {
		service.log.Fatal("gorm.Open error: %s", err)
	}

	db.Create(&entity.Article{Title: keyword.Title, Status: StatusNew})

	/*
		fusionbrainService := NewFusionbrainService(
			cfg,
			service.log,
		)

		result, errGetModels := fusionbrainService.GetModels()
		if errGetModels != nil {
			service.log.Fatal(errGetModels)
		}
		service.log.Info(result.Name, result.Type)

		db, err := gorm.Open(postgres.Open(cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}

		var images []entity.Image
		path := "/home/swenro11/Downloads/"
		db.Where("base64 is not null").Find(&images)
		for _, image := range images {

				errGenerateSlug := fusionbrainService.GenerateSlug(image)
				if errGenerateSlug != nil {
					service.log.Fatal(errGenerateSlug)
				}

			errSaveImage := fusionbrainService.SaveImageToFileSystem(image, path)
			if errSaveImage != nil {
				service.log.Fatal(errSaveImage)
			}
		}

		var img entity.Image
		db.Model(&img).First(&img, "path is not null")
		errDeleteImage := fusionbrainService.DeleteImageFromFileSystem(img)
		if errDeleteImage != nil {
			service.log.Fatal(errDeleteImage)
		}
	*/

	//TODO: create/save link between article & keyword
	//db.Model(&keyword).Update(entity.Keyword{Articles: new []})
	//sdb.Save(&User{Name: "jinzhu", Age: 100})

	//create []Image with generation prompts & StatusNew

	return nil
}
