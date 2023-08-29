package models

type Binding struct {
	BindQueue string
	BindKey   string
}

type Exchange struct {
	Name     string
	Bindings []Binding
}

type RabbitMQNode struct {
	DSN       string
	Exchanges []Exchange
}

type RabbitMQTarget struct {
	DSN          string
	ExchangeName string
	RoutingKey   string
}

type RabbitMQMultiTarget struct {
	DSN          string
	ExchangeName string
	RoutingKeys  []string
}

type RabbitMQSource struct {
	DSN       string
	QueueName string
}

type RabbitMQSourceWithGetback struct {
	DSN                   string
	QueueName             string
	GetbackToExchangeName string
	GetbackKey            string
}
