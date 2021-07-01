package main

import (
	"log"

	rabbitmq "github.com/rsnullptr/go-rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	publisher, returns, err := rabbitmq.NewPublisher(
		"amqp://guest:guest@localhost", amqp.Config{},
		rabbitmq.WithPublisherOptionsLogging,
		func(options *rabbitmq.PublisherOptions) {
			options.BindingExchange = rabbitmq.BindingExchangeOptions{
				Name:          "beautifulExchange1",
				Kind:          "topic",
				Durable:       true,
				AutoDelete:    false,
				Internal:      false,
				NoWait:        true,
				ExchangeArgs:  nil,
				BindingNoWait: false,
			}

			options.RoutingKeys = []string{
				"routing_key1",
			}

			options.QueueDeclare = rabbitmq.QueueDeclareOptions{
				QueueName:       "my_queue1",
				QueueDurable:    true,
				QueueAutoDelete: false,
				QueueExclusive:  false,
				QueueNoWait:     false,
				QueueArgs:       rabbitmq.Table{},
			}

			options.QueueDeclare.QueueArgs["x-queue-type"] = "quorum"
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	err = publisher.Publish(
		[]byte("hello, world1"),
		[]string{"routing_key1"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for r := range returns {
			log.Printf("message returned from server: %s", string(r.Body))
		}
	}()
}
