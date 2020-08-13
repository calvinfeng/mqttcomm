package cmd

import (
	"bytes"
	"crypto/tls"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var Subscriber = &cobra.Command{
	Use:   "subscriber",
	Short: "subscriber subscribes to MQTT broker",
	RunE:  subscribe,
}

func createOnMessage(sig chan os.Signal) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if bytes.Equal(msg.Payload(), []byte("kill")) {
			logrus.Warn("killing process")
			sig <- os.Kill
		} else {
			logrus.Infof("subscriber received message: %s", msg.Payload())
		}
	}
}

func subscribe(cmd *cobra.Command, args []string) error {
	sigKill := make(chan os.Signal)
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("mspurrfection").
		SetCleanSession(true).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})

	opts.OnConnect = func(cli mqtt.Client) {
		token := cli.Subscribe("/purfection", 1, createOnMessage(sigKill))
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
