package cmd

import (
	"context"
	"github.com/calvinfeng/sickmqtt/tdc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math"
	"math/rand"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := cli.SubmitDIOState(tdc.DIOState{
					PinName:             "DIO_A",
					Value:               math.Round(rand.Float64()),
					Direction:           "Output",
					AvailableDirections: []string{"Input", "Output"},
				}); err != nil {
					logrus.WithError(err).Error("failed to submit DIO state to MQTT broker")
				}
			}
		}
	}(ctx)

	<-time.After(1 * time.Hour)
	return nil
}
