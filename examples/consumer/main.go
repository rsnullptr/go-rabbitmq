package main

import (
	"log"

	rabbitmq "github.com/rsnullptr/go-rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	consumer, err := rabbitmq.NewConsumer(
		"amqp://guest:guest@localhost", amqp.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.StartConsuming(
		func(d rabbitmq.Delivery) bool {
			log.Printf("consumed: %v", string(d.Body))
			// true to ACK, false to NACK
			return true
		},
		"my_queue1",
		[]string{"routing_key1", "routing_key_2"},
		rabbitmq.WithConsumeOptionsConcurrency(10),
		func(options *rabbitmq.ConsumeOptions) {
			rabbitmq.WithQueueDeclareOptionsDurable(&options.QueueDeclare)
			rabbitmq.WithQueueDeclareOptionsQuorum(&options.QueueDeclare)
			rabbitmq.WithBindingExchangeOptionsExchangeName("events", &options.BindingExchange)
			rabbitmq.WithBindingExchangeOptionsExchangeKind("topic", &options.BindingExchange)
			rabbitmq.WithBindingExchangeOptionsExchangeDurable(&options.BindingExchange)
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// block main thread so consumers run forever
	forever := make(chan struct{})
	<-forever
}
