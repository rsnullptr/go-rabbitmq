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
		"my_queue7",
		[]string{"publisher-key-1"},
		rabbitmq.WithConsumeOptionsConcurrency(10),
		func(options *rabbitmq.ConsumeOptions) {
			options.BindingExchange = rabbitmq.BindingExchangeOptions{
				DoBinding:     true,
				Name:          "pubexch",
				Kind:          "topic",
				Durable:       true,
				AutoDelete:    false,
				Internal:      false,
				NoWait:        true,
				ExchangeArgs:  nil,
				BindingNoWait: false,
			}

			options.QueueDeclare = rabbitmq.QueueDeclareOptions{
				DoDeclare:       true,
				QueueName:       "my_queue7",
				QueueDurable:    true,
				QueueAutoDelete: false,
				QueueExclusive:  false,
				QueueNoWait:     false,
				QueueArgs:       rabbitmq.Table{},
			}

			options.QueueDeclare.QueueArgs["x-queue-type"] = "quorum"

			options.Qos = rabbitmq.QosOptions{
				QOSPrefetchCount: 1,
				QOSPrefetchSize:  0,
				QOSGlobal:        false,
			}
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// block main thread so consumers run forever
	forever := make(chan struct{})
	<-forever
}
