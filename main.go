package main

import (
	"github.com/calvinfeng/mqttcomm/cmd"
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
	root.AddCommand(cmd.Subscriber, cmd.Publisher)
	if err := root.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
