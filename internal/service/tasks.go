// Package Tasks implements Golang Job Scheduling
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/redis/go-redis/v9"
	"github.com/swenro11/stribog/config"

	log "github.com/swenro11/stribog/pkg/logger"
	"github.com/swenro11/stribog/pkg/postgres"
	"github.com/swenro11/stribog/pkg/rabbitmq/rmq_rpc/client"
)

type TasksService struct {
	repo PoolRepo
	log  *log.Logger
}

func NewTasksService(r PoolRepo, l *log.Logger) *TasksService {
	return &TasksService{
		repo: r,
		log:  l,
	}
}

func (service *TasksService) StartTasks(cfg *config.Config, pg *postgres.Postgres) {
	service.log.Info("StartTasks")

	gocron.Clear()

	//gocron.Every(1).Minute().From(gocron.NextTick()).Do(service.EveryMinuteTask, cfg, pg)
	gocron.Every(10).Minute().From(gocron.NextTick()).Do(service.EveryTenMinuteTask, cfg, pg)
	//gocron.Every(24).Hours().From(gocron.NextTick()).Do(service.EveryDayTask, cfg, pg)

	<-gocron.Start()
}

func (service *TasksService) EveryMinuteTask(cfg *config.Config, pg *postgres.Postgres) {
	service.log.Info("Start EveryMinuteTasks")

	ctx := context.Background()
	tmsp := time.Now().UnixMicro()

	msg := fmt.Sprintf("Start EveryMinuteTasks, time = %s", time.Now())
	service.log.Mongo(ctx, tmsp, msg)

	msg = fmt.Sprintf("End EveryMinuteTasks, check difference between TsUid, time = %s", time.Now())
	service.log.Mongo(ctx, tmsp, msg)
	service.log.Info("End everyMinuteTasks")
}

func (service *TasksService) EveryTenMinuteTask(cfg *config.Config, pg *postgres.Postgres) {
	service.log.Info("Start EveryTenMinuteTask")

	/*
		cohereService := NewCohereService(
			cfg,
			service.log,
		)

		result, errTextGeneration := cohereService.TextGeneration("What is the NATO purpose?")
		if errTextGeneration != nil {
			service.log.Fatal(errTextGeneration)
		}
	*/

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
		articleService := NewArticleService(
			cfg,
			service.log,
		)
		err := articleService.CreateArticle("Test")
		if err != nil {
			service.log.Fatal(err)
		}
	*/
	keywordService := NewKeywordService(
		cfg,
		service.log,
	)
	err := keywordService.CreateKeyword("Test")
	if err != nil {
		service.log.Fatal(err)
	}

	service.log.Info("End everyMinuteTasks")
}

func (service *TasksService) EveryDayTask(cfg *config.Config, pg *postgres.Postgres) {
	service.log.Info("Start EveryDayTask")

	//ctx := context.Background()

	service.log.Info("End EveryDayTask")
}

func (service *TasksService) CheckRabbit(cfg *config.Config, ctx context.Context) {
	//Test RabbitMQ
	rmqClient, err := client.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange)
	if err != nil {
		service.log.Fatal("RabbitMQ RPC Client - init error - client.New")
	}
	defer func() {
		err = rmqClient.Shutdown()
		if err != nil {
			service.log.Fatal("RabbitMQ RPC Client - shutdown error - rmqClient.RemoteCall", err)
		}
	}()
	var answer string

	//TODO: fix
	/*
		"message":"RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall(checkRabbit)%!(EXTRA *fmt.wrapError=rmq_rpc client - Client - RemoteCall - json.Unmarshal: json: cannot unmarshal object into Go value of type string)"}
	*/
	err = rmqClient.RemoteCall("CheckRabbit", nil, &answer)
	if err != nil {
		service.log.Fatal("RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall(checkRabbit)", err)
	}

	if len(answer) > 0 {
		service.log.Info(answer)
	} else {
		service.log.Fatal("CheckRabbit answer is empty")
	}
}

func (service *TasksService) CheckRedis(cfg *config.Config, ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: "",
		DB:       0,
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		service.log.Fatal(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		service.log.Fatal(err)
	}
	service.log.Info("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		service.log.Info("key2 does not exist")
	} else if err != nil {
		service.log.Fatal(err)
	} else {
		service.log.Info("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

// Task execute in RabbitMQ controller
func (service *TasksService) CheckRabbitTask() string {
	msg := fmt.Sprintf("CheckRabbitTask, time = %s", time.Now())
	ctx := context.Background()
	//tmsp нужно вытаскивать или из Redis или из DB, для теста и так сойдет
	tmsp := time.Now().UnixMicro()
	service.log.Mongo(ctx, tmsp, msg)

	return msg
}
