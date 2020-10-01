package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/calvinfeng/sickmqtt/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

type BigLocker struct {
	sync.Mutex
}

func main() {
	l := &BigLocker{
		Mutex: sync.Mutex{},
	}

	l.Lock()
	go func() {
		<-time.After(10 * time.Second)
		l.Unlock()
	}()

	l.Lock()
	fmt.Println("Hello, playground")
	l.Unlock()

	root := &cobra.Command{
		Use: "sickmqtt",
	}
	root.AddCommand(cmd.Subscriber, cmd.Publisher, cmd.RunClient)
	if err := root.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
