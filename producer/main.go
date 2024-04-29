package main

import (
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
	config := newConfig()

	logger.Info("Starting the producer", watermill.LogFields{})
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:     viper.GetStringSlice("Kafka.Brokers"),
			Marshaler:   kafka.DefaultMarshaler{},
			OTELEnabled: true,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	setupWorkers(setupWorkersConfig{
		Count:             config.GetInt("Workers"),
		MessagesPerSecond: config.GetInt("MessagesPerSecond"),
		Publisher:         publisher,
	})

	logger.Info("All messages published", nil)
}

type setupWorkersConfig struct {
	Count             int
	MessagesPerSecond int
	Publisher         *kafka.Publisher
}

func setupWorkers(config setupWorkersConfig) {
	closeCh := make(chan struct{})
	workersGroup := &sync.WaitGroup{}
	workersGroup.Add(config.Count)

	for i := 0; i < config.Count; i++ {
		go worker.Create(config.MessagesPerSecond, config.Publisher, workersGroup, closeCh)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	close(closeCh)

	workersGroup.Wait()
}

func newConfig() *viper.Viper {
	config := viper.New()
	config.AddConfigPath(".")
	config.AutomaticEnv()
	config.SetConfigName("settings")
	config.SetConfigType("json")

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return config
}
