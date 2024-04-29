package router

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/dprcoles/kafka-pubsub/consumer/handlers"
	"github.com/spf13/viper"
	"time"
)

func Create(config *viper.Viper, pub *kafka.Publisher, marshaler kafka.Unmarshaler, logger watermill.LoggerAdapter) (*message.Router, error) {
	r, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	retryMiddleware := middleware.Retry{
		MaxRetries:      1,
		InitialInterval: time.Millisecond * 10,
	}

	poisonQueue, err := middleware.PoisonQueue(pub, config.GetString("Kafka.Topics.PoisonQueue"))
	if err != nil {
		panic(err)
	}

	r.AddMiddleware(
		// Recover from panics in handlers
		middleware.Recoverer,
		// Limit incoming messages to 10 per second
		middleware.NewThrottle(10, time.Second).Middleware,
		// If the retries limit is exceeded (see retryMiddleware below), the message is sent
		// to the poison queue (published to poison_queue topic)
		poisonQueue,
		// retries message processing if an error occurred in the handler
		retryMiddleware.Middleware,
		// Correlation ID middleware adds the correlation ID of the consumed message to each produced message.
		middleware.CorrelationID,
		// Simulate errors or panics from handler
		middleware.RandomFail(0.01),
		middleware.RandomPanic(0.01),
	)

	// Close the router when a SIGTERM is sent
	r.AddPlugin(plugin.SignalsHandler)

	// Handler that counts consumed posts
	r.AddHandler(
		"posts_counter",
		config.GetString("Kafka.Topics.PostsPublished"),
		createSubscriber(config, "posts_counter", marshaler, logger),
		config.GetString("Kafka.Topics.PostsCounted"),
		pub,
		handlers.NewPostsCounter().Count,
	)

	// Handler that generates "feed" from consumed posts
	r.AddNoPublisherHandler(
		"feed_generator",
		config.GetString("Kafka.Topics.PostsPublished"),
		createSubscriber(config, "feed_generator", marshaler, logger),
		handlers.NewFeedGenerator().UpdateFeed,
	)

	return r, nil
}

func createSubscriber(config *viper.Viper, consumerGroup string, marshaler kafka.Unmarshaler, logger watermill.LoggerAdapter) message.Subscriber {
	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       config.GetStringSlice("Kafka.Brokers"),
			Unmarshaler:   marshaler,
			ConsumerGroup: consumerGroup,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	return sub
}
