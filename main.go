package main

import (
	"github.com/calvinfeng/sickmqtt/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	root := &cobra.Command{
		Use: "mqttcomm",
	}
	root.AddCommand(cmd.Subscriber, cmd.Publisher, cmd.RunClient)
	if err := root.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
