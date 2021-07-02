package main

import (
	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	"log"
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
		"my_queue5",
		[]string{"routing_key_7"},
		rabbitmq.WithConsumeOptionsConcurrency(10),
		func(options *rabbitmq.ConsumeOptions) {
			rabbitmq.WithQueueDeclare(&options.QueueDeclare)
			rabbitmq.WithQueueDeclareOptionsDurable(&options.QueueDeclare)
			rabbitmq.WithQueueDeclareOptionsQuorum(&options.QueueDeclare)
			rabbitmq.WithBindingExchangeOptionsExchangeName("events", &options.BindingExchange)
			rabbitmq.WithBindingExchangeOptionsExchangeKind("topic", &options.BindingExchange)
			rabbitmq.WithBindingExchangeOptionsExchangeDurable(&options.BindingExchange)
			options.Qos.QOSPrefetchCount = 1
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// block main thread so consumers run forever
	forever := make(chan struct{})
	<-forever
}
