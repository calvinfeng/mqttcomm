package tdc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
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
	LEDStatesAll   = "/tdce/led/#"
	LEDStatesRead  = "/tdce/led/out"
	LEDStatesWrite = "/tdce/led/in"
	DIOStatesRead  = "/tdce/dio/state/out"
	DIOStatesWrite = "/tdce/dio/state/in"
	AINValueRead   = "/tdce/analog-inputs/value"
	AINModeRead    = "/tdce/analog-inputs/mode"
	AINStateRead   = "/tdce/analog-inputs/state"
	ErrorRead      = "/tdce/error"
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
			ErrorRead:      1,
		}, func(_ mqtt.Client, msg mqtt.Message) {
			go m.handleMQTTMessage(msg)
		})
	}

	cli := mqtt.NewClient(opts)
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.cli = cli
	return nil
}

func (m *mqttclient) SubmitDIOState(state DIOState) error {
	if m.cli == nil || !m.cli.IsConnected() {
		return errors.New("MQTT client is not connected to broker, please call Connect() first")
	}

	states := []DIOState{state}

	jsonPayload, err := json.Marshal(states)
	if err != nil {
		return err
	}

	// Retained is set to true because we want to make best effort of delivering message to a TDC.
	if token := m.cli.Publish(DIOStatesWrite, 1, true, jsonPayload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *mqttclient) handleMQTTMessage(msg mqtt.Message) {
	switch msg.Topic() {
	case DIOStatesRead:
		logrus.Infof("topic %s", msg.Topic())
		var dioState DIOState
		if err := json.Unmarshal(msg.Payload(), &dioState); err != nil {
			logrus.WithError(err).Error("failed to parse payload %s", string(msg.Payload()))
		} else {
			logrus.Infof("successfully parsed message payload %v", dioState)
		}

	case AINStateRead:
		logrus.Infof("topic %s", msg.Topic())
		var stateChange AINStateChange
		if err := json.Unmarshal(msg.Payload(), &stateChange); err != nil {
			logrus.WithError(err).Error("failed to parse payload %s", string(msg.Payload()))
		} else {
			logrus.Infof("successfully parsed message payload %v", stateChange)
		}

	case AINValueRead:
		logrus.Infof("topic %s", msg.Topic())
		var valueChange AINValueChange
		if err := json.Unmarshal(msg.Payload(), &valueChange); err != nil {
			logrus.WithError(err).Error("failed to parse payload %s", string(msg.Payload()))
		} else {
			logrus.Infof("successfully parsed message payload %v", valueChange)
		}

	case AINModeRead:
		logrus.Infof("topic %s", msg.Topic())
		var valueChange AINModeChange
		if err := json.Unmarshal(msg.Payload(), &valueChange); err != nil {
			logrus.WithError(err).Error("failed to parse payload %s", string(msg.Payload()))
		} else {
			logrus.Infof("successfully parsed message payload %v", valueChange)
		}

	case ErrorRead:
		logrus.Infof("topic %s", msg.Topic())
		logrus.Infof("error %s", msg.Payload())

	case LEDStatesRead:
		logrus.Infof("topic %s", msg.Topic())
		var ledStates []LEDState
		if err := json.Unmarshal(msg.Payload(), &ledStates); err != nil {
			logrus.WithError(err).Error("failed to parse payload %s", string(msg.Payload()))
		} else {
			logrus.Infof("successfully parsed message payload %v", ledStates)
		}
	}
}
