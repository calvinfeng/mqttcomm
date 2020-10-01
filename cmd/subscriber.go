package cmd

import (
	"bytes"
	"crypto/tls"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Subscriber = &cobra.Command{
	Use:   "subscriber",
	Short: "subscriber subscribes to MQTT broker",
	RunE:  subscribe,
}

const readDIOState = "/tdc/+/dio/out"

func createOnMessage(sig chan os.Signal) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if bytes.Equal(msg.Payload(), []byte("kill")) {
			logrus.Warn("killing process")
			sig <- os.Kill
		} else {
			if msg.Topic()[0] != '/' {
				logrus.Warnf("topic %s has invalid format, it does not start with /", msg.Topic())
				return
			}

			parts := strings.Split(msg.Topic(), "/")
			if len(parts) < 3 {
				logrus.Warnf("%s should have more than 3 parts if split by /", msg.Topic())
				return
			}

			deviceID := parts[2]
			switch strings.Replace(msg.Topic(), deviceID, "+", 1) {
			case readDIOState:
				logrus.Infof("subscriber received topic: %s", msg.Topic())
				logrus.Infof("subscriber received message: %s", msg.Payload())
			default:
				logrus.Errorf("topic %s is unrecognized", msg.Topic())
			}
		}
	}
}

func subscribe(cmd *cobra.Command, args []string) error {
	sigKill := make(chan os.Signal)
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("osxsubscriber").
		SetCleanSession(true).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})

	opts.OnConnect = func(cli mqtt.Client) {
		token := cli.Subscribe("/tdc/+/dio/out", 0, createOnMessage(sigKill))
		if token.Wait() && token.Error() != nil {
			logrus.Fatal(token.Error())
		}
	}

	subscriber := mqtt.NewClient(opts)
	logrus.Infof("subscriber is connecting to %s", broker)
	token := subscriber.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	logrus.Infof("subscriber is connected to %s", broker)

	<-sigKill
	return nil
}
