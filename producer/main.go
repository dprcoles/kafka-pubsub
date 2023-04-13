package main

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/dprcoles/kafka-pubsub/producer/worker"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	handleConfig()

	logger.Info("Starting the producer", watermill.LogFields{})

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   viper.GetStringSlice("Kafka.Brokers"),
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	setupWorkers(publisher)

	logger.Info("All messages published", nil)
}

func setupWorkers(publisher *kafka.Publisher) {
	workers := viper.GetInt("Workers")

	closeCh := make(chan struct{})
	workersGroup := &sync.WaitGroup{}
	workersGroup.Add(workers)

	for i := 0; i < workers; i++ {
		go worker.Create(publisher, workersGroup, closeCh)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	close(closeCh)

	workersGroup.Wait()
}

func handleConfig() {
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigName("settings")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("Failed to load configuration: %v", err)
	}
}
