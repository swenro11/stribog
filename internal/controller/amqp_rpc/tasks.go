package amqprpc

import (
	"github.com/streadway/amqp"

	"github.com/swenro11/stribog/internal/entity"
	"github.com/swenro11/stribog/internal/service"
	"github.com/swenro11/stribog/pkg/rabbitmq/rmq_rpc/server"
)

type tasksRoutes struct {
	tasks service.Tasks
}

func newTasksRoutes(routes map[string]server.CallHandler, t service.Tasks) {
	r := &tasksRoutes{t}
	{
		//routes["checkProfit"] = r.checkProfit()
		routes["CheckRabbitTask"] = r.CheckRabbitTask()
	}
}

type poolResponse struct {
	Pool []entity.Pool `json:"pool"`
}

type rabbitResponse struct {
	answer string `json:"answer"`
}

// not ended, obsolete
/*
func (r *tasksRoutes) checkProfit() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		pools, err := r.tasksService.CheckProfit(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - tasksRoutes - checkProfit - r.tasksService.CheckProfit: %w", err)
		}
		response := poolResponse{pools}

		return response, nil
	}
}
*/

// simple RabbitMQ task, with logging
func (r *tasksRoutes) CheckRabbitTask() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		answer := r.tasks.CheckRabbitTask()

		response := rabbitResponse{answer}

		return response, nil
	}
}
