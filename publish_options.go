package rabbitmq

// PublisherOptions are used to describe a publisher's configuration.
// Logging set to true will enable the consumer to print to stdout
type PublisherOptions struct {
	Logging         bool
	Logger          Logger
	RoutingKeys     []string
	BindingExchange BindingExchangeOptions
	QueueDeclare    QueueDeclareOptions
	Qos             QosOptions
}

// WithPublisherOptionsLogging sets logging to true on the consumer options
func WithPublisherOptionsLogging(options *PublisherOptions) {
	options.Logging = true
	options.Logger = &stdLogger{}
}

// WithPublisherOptionsLogger sets logging to a custom interface.
// Use WithPublisherOptionsLogging to just log to stdout.
func WithPublisherOptionsLogger(log Logger) func(options *PublisherOptions) {
	return func(options *PublisherOptions) {
		options.Logging = true
		options.Logger = log
	}
}

// PublishOptions are used to control how data is published
type PublishOptions struct {
	Exchange string
	// Mandatory fails to publish if there are no queues
	// bound to the routing key
	Mandatory bool
	// Immediate fails to publish if there are no consumers
	// that can ack bound to the queue on the routing key
	Immediate   bool
	ContentType string
	// Transient or Persistent
	DeliveryMode uint8
	// Expiration time in ms that a message will expire from a queue.
	// See https://www.rabbitmq.com/ttl.html#per-message-ttl-in-publishers
	Expiration string
	Headers    Table
}

// WithPublishOptionsExchange returns a function that sets the exchange to publish to
func WithPublishOptionsExchange(exchange string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Exchange = exchange
	}
}

// WithPublishOptionsMandatory makes the publishing mandatory, which means when a queue is not
// bound to the routing key a message will be sent back on the returns channel for you to handle
func WithPublishOptionsMandatory(options *PublishOptions) {
	options.Mandatory = true
}

// WithPublishOptionsImmediate makes the publishing immediate, which means when a consumer is not available
// to immediately handle the new message, a message will be sent back on the returns channel for you to handle
func WithPublishOptionsImmediate(options *PublishOptions) {
	options.Immediate = true
}

// WithPublishOptionsContentType returns a function that sets the content type, i.e. "application/json"
func WithPublishOptionsContentType(contentType string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.ContentType = contentType
	}
}

// WithPublishOptionsPersistentDelivery sets the message to persist. Transient messages will
// not be restored to durable queues, persistent messages will be restored to
// durable queues and lost on non-durable queues during server restart. By default publishings
// are transient
func WithPublishOptionsPersistentDelivery(options *PublishOptions) {
	options.DeliveryMode = Persistent
}

// WithPublishOptionsExpiration returns a function that sets the expiry/TTL of a message. As per RabbitMq spec, it must be a
// string value in milliseconds.
func WithPublishOptionsExpiration(expiration string) func(options *PublishOptions) {
	return func(options *PublishOptions) {
		options.Expiration = expiration
	}
}

// WithPublishOptionsHeaders returns a function that sets message header values, i.e. "msg-id"
func WithPublishOptionsHeaders(headers Table) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Headers = headers
	}
}
