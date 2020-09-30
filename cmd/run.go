package cmd

import (
	"bufio"
	"fmt"
	"github.com/calvinfeng/sickmqtt/tdc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
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

	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter a state 0 or 1:")
		line, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		message := strings.TrimSuffix(line, "\n")

		if message == "0" {
			if err := cli.SubmitDIOState(tdc.DIOState{
				PinName:             "DIO_A",
				Value:               0,
				Direction:           "Output",
				AvailableDirections: []string{"Input", "Output"},
			}); err != nil {
				logrus.WithError(err).Error("failed to submit DIO state to MQTT broker")
			}
		} else if message == "1"{
			if err := cli.SubmitDIOState(tdc.DIOState{
				PinName:             "DIO_A",
				Value:               1,
				Direction:           "Output",
				AvailableDirections: []string{"Input", "Output"},
			}); err != nil {
				logrus.WithError(err).Error("failed to submit DIO state to MQTT broker")
			}
		} else {
			logrus.Errorf("%s is not a recognized command", message)
		}
	}
	return nil
}
