package cmd

import (
	"bufio"
	"crypto/tls"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var broker = "tcp://test.mosquitto.org:1883"

var Publisher = &cobra.Command{
	Use:   "publisher",
	Short: "publisher publishes messages from stdin to broker",
	RunE:  publish,
}

func publish(cmd *cobra.Command, args []string) error {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("sirpurrfection").
		SetCleanSession(true).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})

	publisher := mqtt.NewClient(opts)
	logrus.Infof("publisher is connecting to %s", broker)
	if token := publisher.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	logrus.Infof("publisher is connected to %s", broker)

	stdin := bufio.NewReader(os.Stdin)
	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		// QoS stands for quality of service
		// 0 for Once
		// 1 for At Least Once (may be delivered more than once)
		logrus.Info("publishing message...")
		token := publisher.Publish("/purfection", 1, false, message)
		if token.Wait() && token.Error() != nil {
			logrus.Error(token.Error())
		} else {
			logrus.Info("message is successfully published")
		}
	}
	return nil
}
