// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/swenro11/stribog/config"
	amqprpc "github.com/swenro11/stribog/internal/controller/amqp_rpc"
	"github.com/swenro11/stribog/internal/service"
	"github.com/swenro11/stribog/internal/service/repo"
	"github.com/swenro11/stribog/pkg/logger"
	"github.com/swenro11/stribog/pkg/postgres"
	"github.com/swenro11/stribog/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(context.Background(), cfg.Log.Level, cfg.PARAM.TgBotApi, cfg.PARAM.TgChatId, cfg.Mongo.URI, cfg.Mongo.DB)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// RabbitMQ RPC Server
	/*
		translationService := service.NewTranslationService(
			repo.NewTranslationRepo(pg),
			webapi.New(),
		)

		rmqRouter := amqprpc.NewTranslationRouter(translationService)

		rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, log)
		if err != nil {
			log.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
		}
	*/

	tasksService := service.NewTasksService(repo.NewPoolRepo(pg), log)

	rmqRouter := amqprpc.NewTasksRouter(tasksService)

	rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, log)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	/*
		handler := gin.New()
		v1.NewRouter(handler, l, translationService)
		httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	*/

	// Tasks
	tasksService.StartTasks(cfg, pg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	/*
		case err = <-httpServer.Notify():
			l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	*/
	case err = <-rmqServer.Notify():
		log.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	/*
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	*/

	err = rmqServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

}
