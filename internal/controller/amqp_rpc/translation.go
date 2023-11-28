package amqprpc

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/swenro11/stribog/internal/entity"
	"github.com/swenro11/stribog/internal/service"
	"github.com/swenro11/stribog/pkg/rabbitmq/rmq_rpc/server"
)

type translationRoutes struct {
	translationService service.Translation
}

func newTranslationRoutes(routes map[string]server.CallHandler, t service.Translation) {
	r := &translationRoutes{t}
	{
		routes["getHistory"] = r.getHistory()
	}
}

type historyResponse struct {
	History []entity.Translation `json:"history"`
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		translations, err := r.translationService.History(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationService.History: %w", err)
		}

		response := historyResponse{translations}

		return response, nil
	}
}
