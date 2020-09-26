package tdc

type MQTTClientOptions struct {
	BrokerAddr string
}

type MQTTClient interface {
	Connect() error
}
