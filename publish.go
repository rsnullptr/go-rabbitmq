package rabbitmq

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

// DeliveryMode. Transient means higher throughput but messages will not be
// restored on broker restart. The delivery mode of publishings is unrelated
// to the durability of the queues they reside on. Transient messages will
// not be restored to durable queues, persistent messages will be restored to
// durable queues and lost on non-durable queues during server restart.
//
// This remains typed as uint8 to match Publishing.DeliveryMode. Other
// delivery modes specific to custom queue implementations are not enumerated
// here.
const (
	Transient  uint8 = amqp.Transient
	Persistent uint8 = amqp.Persistent
)

// Return captures a flattened struct of fields returned by the server when a
// Publishing is unable to be delivered either due to the `mandatory` flag set
// and no route found, or `immediate` flag set and no free consumer.
type Return struct {
	amqp.Return
}

// Publisher allows you to publish messages safely across an open connection
type Publisher struct {
	chManager *channelManager

	notifyFlowChan chan bool

	disablePublishDueToFlow    bool
	disablePublishDueToFlowMux *sync.RWMutex

	logger Logger
}

// NewPublisher returns a new publisher with an open channel to the cluster.
// If you plan to enforce mandatory or immediate publishing, those failures will be reported
// on the channel of Returns that you should setup a listener on.
// Flow controls are automatically handled as they are sent from the server, and publishing
// will fail with an error when the server is requesting a slowdown
func NewPublisher(url string, config amqp.Config, optionFuncs ...func(*PublisherOptions)) (Publisher, <-chan Return, error) {
	options := &PublisherOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if options.Logger == nil {
		options.Logger = &noLogger{} // default no logging
	}

	chManager, err := newChannelManager(url, config, options.Logger)
	if err != nil {
		return Publisher{}, nil, err
	}

	publisher := Publisher{
		chManager:                  chManager,
		notifyFlowChan:             make(chan bool),
		disablePublishDueToFlow:    false,
		disablePublishDueToFlowMux: &sync.RWMutex{},
		logger:                     options.Logger,
	}

	returnAMQPChan := make(chan amqp.Return)
	returnChan := make(chan Return)
	returnAMQPChan = publisher.chManager.channel.NotifyReturn(returnAMQPChan)
	go func() {
		for ret := range returnAMQPChan {
			returnChan <- Return{
				ret,
			}
		}
	}()

	publisher.notifyFlowChan = publisher.chManager.channel.NotifyFlow(publisher.notifyFlowChan)
	go publisher.startNotifyFlowHandler()

	publisher.chManager.channelMux.RLock()
	defer publisher.chManager.channelMux.RUnlock()

	// valid queue name, declare it
	if options.QueueDeclare.DoDeclare {

		_, err = publisher.chManager.channel.QueueDeclare(
			options.QueueDeclare.QueueName,
			options.QueueDeclare.QueueDurable,
			options.QueueDeclare.QueueAutoDelete,
			options.QueueDeclare.QueueExclusive,
			options.QueueDeclare.QueueNoWait,
			tableToAMQPTable(options.QueueDeclare.QueueArgs),
		)

		if err != nil {
			return publisher, returnChan, err
		}
	}

	//valid name, bind it
	if options.BindingExchange.DoBinding {
		exchange := options.BindingExchange

		err = publisher.chManager.channel.ExchangeDeclare(
			exchange.Name,
			exchange.Kind,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			tableToAMQPTable(exchange.ExchangeArgs),
		)

		if err != nil {
			return publisher, returnChan, err
		}

		for _, routingKey := range options.RoutingKeys {
			err = publisher.chManager.channel.QueueBind(
				options.QueueDeclare.QueueName,
				routingKey,
				exchange.Name,
				options.BindingExchange.BindingNoWait,
				tableToAMQPTable(options.BindingExchange.BindingArgs),
			)
			if err != nil {
				return publisher, returnChan, err
			}
		}
	}

	err = publisher.chManager.channel.Qos(
		options.Qos.QOSPrefetchCount,
		options.Qos.QOSPrefetchSize,
		options.Qos.QOSGlobal,
	)

	return publisher, returnChan, err
}

// Publish publishes the provided data to the given routing keys over the connection
func (publisher *Publisher) Publish(
	data []byte,
	routingKeys []string,
	optionFuncs ...func(*PublishOptions),
) error {
	publisher.disablePublishDueToFlowMux.RLock()
	if publisher.disablePublishDueToFlow {
		return fmt.Errorf("publishing blocked due to high flow on the server")
	}
	publisher.disablePublishDueToFlowMux.RUnlock()

	options := &PublishOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if options.DeliveryMode == 0 {
		options.DeliveryMode = Transient
	}

	for _, routingKey := range routingKeys {
		var message = amqp.Publishing{}
		message.ContentType = options.ContentType
		message.DeliveryMode = options.DeliveryMode
		message.Body = data
		message.Headers = tableToAMQPTable(options.Headers)
		message.Expiration = options.Expiration

		// Actual publish.
		err := publisher.chManager.channel.Publish(
			options.Exchange,
			routingKey,
			options.Mandatory,
			options.Immediate,
			message,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// StopPublishing stops the publishing of messages.
// The publisher should be discarded as it's not safe for re-use
func (publisher Publisher) StopPublishing() {
	publisher.chManager.channel.Close()
	publisher.chManager.connection.Close()
}

func (publisher *Publisher) startNotifyFlowHandler() {
	// Listeners for active=true flow control.  When true is sent to a listener,
	// publishing should pause until false is sent to listeners.
	for ok := range publisher.notifyFlowChan {
		publisher.disablePublishDueToFlowMux.Lock()
		if ok {
			publisher.logger.Printf("pausing publishing due to flow request from server")
			publisher.disablePublishDueToFlow = true
		} else {
			publisher.disablePublishDueToFlow = false
			publisher.logger.Printf("resuming publishing due to flow request from server")
		}
		publisher.disablePublishDueToFlowMux.Unlock()
	}
}
