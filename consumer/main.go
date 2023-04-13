package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/dprcoles/kafka-pubsub/consumer/router"
	"github.com/spf13/viper"
)

func main() {
	marshaler := kafka.DefaultMarshaler{}

	logger := watermill.NewStdLogger(false, false)
	logger.Info("Starting the consumer", nil)

	handleConfig()

	pub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   viper.GetStringSlice("Kafka.Brokers"),
			Marshaler: marshaler,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	r, err := router.Create(pub, marshaler, logger)

	if err = r.Run(context.Background()); err != nil {
		panic(err)
	}
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
