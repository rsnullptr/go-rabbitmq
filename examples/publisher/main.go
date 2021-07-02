package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
)

func main() {
	do()
}

func do() {
	publisher, returns, err := rabbitmq.NewPublisher(
		"amqp://guest:guest@localhost", amqp.Config{},
		rabbitmq.WithPublisherOptionsLogging,
		func(options *rabbitmq.PublisherOptions) {
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

			options.RoutingKeys = []string{
				"publisher-key-1",
			}

			options.QueueDeclare = rabbitmq.QueueDeclareOptions{
				DoDeclare:       true,
				QueueName:       "pub_queue",
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
				QOSGlobal:        true,
			}
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for r := range returns {
			log.Printf("message returned from server: %s", string(r.Body))
		}
	}()

	idx := int64(0)
	for {
		err = publisher.Publish(
			[]byte(fmt.Sprintf("AAA %d", idx)),
			[]string{"publisher-key-1"},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsMandatory,
			rabbitmq.WithPublishOptionsPersistentDelivery,
			rabbitmq.WithPublishOptionsExchange("pubexch"),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(os.Stderr, "%d \n", idx)
		idx++
		time.Sleep(time.Millisecond * 1000)
	}
}
