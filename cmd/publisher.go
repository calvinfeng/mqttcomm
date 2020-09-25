package cmd

import (
	"bufio"
	"crypto/tls"
	"io"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Mosquitto
// var broker = "10.100.3.14:1883"

// Linux
var broker = "10.100.3.53:1883"

// Mac
// var broker = "192.168.100.76:1883"

var Publisher = &cobra.Command{
	Use:   "publisher",
	Short: "publisher publishes messages from stdin to broker",
	RunE:  publish,
}

func publish(cmd *cobra.Command, args []string) error {
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	//tlsConfig, err := security.NewTLSConfig()
	//if err != nil {
	//	return err
	//}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("osxpublisher").
		SetCleanSession(true).
		SetTLSConfig(tlsConfig)

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
		token := publisher.Publish("/input/totdc", 1, false, message)
		if token.Wait() && token.Error() != nil {
			logrus.Error(token.Error())
		} else {
			logrus.Info("message is successfully published")
		}
	}
	return nil
}
