package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NYTimes/gizmo/pubsub"
	"github.com/Sirupsen/logrus"
	"github.com/markgunnels/gizmo/pubsub/amqp"
)

var (
	Log = logrus.New()
	sub pubsub.Subscriber
)

func Init() {
	cfg := amqp.LoadAMQPConfigFromEnv()

	Log.Out = os.Stderr
	pubsub.Log = Log

	var err error

	sub, err = amqp.NewAMQPSubscriber(cfg)
	if err != nil {
		Log.Fatal("unable to init AMQP: ", err)
	}
}

func Run() (err error) {
	stream := sub.Start()

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		Log.Infof("received kill signal %s", <-ch)
		err = sub.Stop()
	}()

	for msg := range stream {
		message := string(msg.Message())

		fmt.Println(message)
		if err = msg.Done(); err != nil {
			Log.WithFields(logrus.Fields{
				"message": message,
			}).Error("unable to delete message from AMQP: ", err)
			return err
		}
	}

	return err
}
