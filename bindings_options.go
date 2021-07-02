package rabbitmq

// BindingExchangeOptions are used when binding to an exchange.
// it will verify the exchange is created before binding to it.
type BindingExchangeOptions struct {
	DoBinding     bool
	Name          string
	Kind          string
	Durable       bool
	AutoDelete    bool
	Internal      bool
	NoWait        bool
	BindingNoWait bool
	BindingArgs   Table
	ExchangeArgs  Table
}

// QueueDeclareOptions arguments to declare a queue
type QueueDeclareOptions struct {
	DoDeclare       bool
	QueueName       string
	QueueDurable    bool
	QueueAutoDelete bool
	QueueExclusive  bool
	QueueNoWait     bool
	QueueArgs       Table
}

// QosOptions configuration
type QosOptions struct {
	QOSPrefetchCount int
	QOSPrefetchSize  int
	QOSGlobal        bool
}

// WithQueueDeclare define if queue will be declare or not
func WithQueueDeclare(options *QueueDeclareOptions) {
	options.DoDeclare = true
}

// WithQueueDeclareOptionsDurable sets the queue to durable, which means it won't
// be destroyed when the server restarts. It must only be bound to durable exchanges
func WithQueueDeclareOptionsDurable(options *QueueDeclareOptions) {
	options.QueueDurable = true
}

// WithQueueDeclareOptionsAutoDelete sets the queue to auto delete, which means it will
// be deleted when there are no more consumers on it
func WithQueueDeclareOptionsAutoDelete(options *QueueDeclareOptions) {
	options.QueueAutoDelete = true
}

// WithQueueDeclareOptionsExclusive sets the queue to exclusive, which means
// it's are only accessible by the connection that declares it and
// will be deleted when the connection closes. Channels on other connections
// will receive an error when attempting to declare, bind, consume, purge or
// delete a queue with the same name.
func WithQueueDeclareOptionsExclusive(options *QueueDeclareOptions) {
	options.QueueExclusive = true
}

// WithQueueDeclareOptionsNoWait sets the queue to nowait, which means
// the queue will assume to be declared on the server.  A
// channel exception will arrive if the conditions are met for existing queues
// or attempting to modify an existing queue from a different connection.
func WithQueueDeclareOptionsNoWait(options *QueueDeclareOptions) {
	options.QueueNoWait = true
}

// WithQueueDeclareOptionsQuorum sets the queue a quorum type, which means multiple nodes
// in the cluster will have the messages distributed amongst them for higher reliability
func WithQueueDeclareOptionsQuorum(options *QueueDeclareOptions) {
	if options.QueueArgs == nil {
		options.QueueArgs = Table{}
	}
	options.QueueArgs["x-queue-type"] = "quorum"
}

// getBindingExchangeOptionsOrSetDefault returns pointer to current BindingExchange options. if no BindingExchange options are set yet, it will set it with default values.
func getBindingExchangeOptionsOrSetDefault(options *BindingExchangeOptions) *BindingExchangeOptions {
	if options == nil {
		options = &BindingExchangeOptions{
			Name:         "",
			Kind:         "direct",
			Durable:      false,
			AutoDelete:   false,
			Internal:     false,
			NoWait:       false,
			ExchangeArgs: nil,
		}
	}
	return options
}

// WithBindingExchangeOptionsExchangeName returns a function that sets the exchange name the queue will be bound to
func WithBindingExchangeOptionsExchangeName(name string, options *BindingExchangeOptions) {
	options.Name = name
}

// WithBindingExchangeOptionsExchangeKind returns a function that sets the binding exchange kind/type
func WithBindingExchangeOptionsExchangeKind(kind string, options *BindingExchangeOptions) {
	options.Kind = kind
}

// WithBindingExchangeOptionsExchangeDurable returns a function that sets the binding exchange durable flag
func WithBindingExchangeOptionsExchangeDurable(options *BindingExchangeOptions) {
	options.Durable = true
}

// WithBindingExchangeOptionsExchangeAutoDelete returns a function that sets the binding exchange autoDelete flag
func WithBindingExchangeOptionsExchangeAutoDelete(options *BindingExchangeOptions) {
	options.AutoDelete = true
}

// WithBindingExchangeOptionsExchangeInternal returns a function that sets the binding exchange internal flag
func WithBindingExchangeOptionsExchangeInternal(options *BindingExchangeOptions) {
	options.Internal = true
}

// WithBindingExchangeOptionsExchangeNoWait returns a function that sets the binding exchange noWait flag
func WithBindingExchangeOptionsExchangeNoWait(options *BindingExchangeOptions) {
	options.NoWait = true
}

// WithBindingExchangeOptionsExchangeArgs returns a function that sets the binding exchange arguments that are specific to the server's implementation of the exchange
func WithBindingExchangeOptionsExchangeArgs(args Table, options *BindingExchangeOptions) {
	options.ExchangeArgs = args
}

// WithBindingExchangeOptionsNoWait sets the bindings to nowait, which means if the queue can not be bound
// the channel will not be closed with an error.
func WithBindingExchangeOptionsNoWait(options *BindingExchangeOptions) {
	options.BindingNoWait = true
}

// WithBindingExchange define if the exchange is to be binded or not
func WithBindingExchange(options *BindingExchangeOptions) {
	options.DoBinding = true
}
