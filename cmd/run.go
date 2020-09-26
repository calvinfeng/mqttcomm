package cmd

import (
	"github.com/calvinfeng/sickmqtt/tdc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var RunClient = &cobra.Command{
	Use:   "run_client",
	Short: "run client to talk to SICK TDC via MQTT broker",
	RunE:  runClient,
}

func runClient(cmd *cobra.Command, args []string) error {
	cli := tdc.NewMQTTClient(tdc.MQTTClientOptions{BrokerAddr: "10.100.3.57:1883"})
	if err := cli.Connect(); err != nil {
		return err
	}
	logrus.Infof("MQTT client is connected")
	<-time.After(1 * time.Hour)
	return nil
}
