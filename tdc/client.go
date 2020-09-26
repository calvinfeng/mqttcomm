package tdc

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(opts MQTTClientOptions) MQTTClient {
	return &mqttclient{
		opts: opts,
	}
}

type mqttclient struct {
	opts MQTTClientOptions
	cli  mqtt.Client
}

const (
	LEDStatesRead  = "/tdce/led/out"
	LEDStatesWrite = "/tdce/led/in"
	DIOStatesRead  = "/tdce/dio/state/out"
	DIOStatesWrite = "/tdce/dio/state/in"
	AINValueRead   = "/tdce/analog-inputs/value"
	AINModeRead    = "/tdce/analog-inputs/mode"
	AINStateRead   = "/tdce/analog-inputs/state"
)

func (m *mqttclient) Connect() error {
	opts := mqtt.NewClientOptions().
		AddBroker(m.opts.BrokerAddr).
		SetClientID("random").
		SetCleanSession(true).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})

	opts.OnConnect = func(cli mqtt.Client) {
		cli.SubscribeMultiple(map[string]byte{
			LEDStatesRead:  1,
			LEDStatesWrite: 1,
			DIOStatesRead:  1,
			DIOStatesWrite: 1,
			AINValueRead:   1,
			AINModeRead:    1,
			AINStateRead:   1,
		}, func(_ mqtt.Client, msg mqtt.Message) {
			fmt.Println(msg.Topic())
			fmt.Println(string(msg.Payload()))
		})
	}

	cli := mqtt.NewClient(opts)
	if token := cli.Connect(); token.Wait() || token.Error() != nil {
		return token.Error()
	}
	m.cli = cli
	return nil
}
