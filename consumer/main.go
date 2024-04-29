package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/dprcoles/kafka-pubsub/consumer/router"
	"github.com/spf13/viper"
)

func main() {
	marshaler := kafka.DefaultMarshaler{}

	logger := watermill.NewStdLogger(false, false)
	logger.Info("Starting the consumer", nil)

	config := newConfig()

	pub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:     config.GetStringSlice("Kafka.Brokers"),
			Marshaler:   marshaler,
			OTELEnabled: true,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	r, err := router.Create(config, pub, marshaler, logger)
	if err != nil {
		panic(err)
	}

	if err = r.Run(context.Background()); err != nil {
		panic(err)
	}
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
