package main

import (
	"log/slog"
	"os"
	"time"

	ssepubsub "github.com/Christian-007/fit-forge-notification/internal/app/sse/delivery/pubsub"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/appdependency"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/decorator"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func NewWatermillRouter(watermillLogger watermill.LoggerAdapter, appDependencies appdependency.AppDependency) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		appDependencies.Logger.Error("failed to create Watermill Router",
			slog.String("error", err.Error()),
		)
		panic(err)
	}

	// Exponential backoff
	router.AddMiddleware(middleware.Retry{
		MaxRetries:      5,
		InitialInterval: time.Millisecond * 500,
		MaxInterval:     time.Second * 30,
		Multiplier:      2,
		Logger:          watermillLogger,
	}.Middleware)

	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			correlationId := msg.Metadata.Get(decorator.CorrelationIdMetadataKey)
			appDependencies.Logger.Info("handling a message",
				slog.String("correlation_id", correlationId),
				slog.String("message_id", msg.UUID),
				slog.Any("metadata", msg.Metadata),
				slog.String("payload", string(msg.Payload)),
			)
			return next(msg)
		}
	})

	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			ProjectID: os.Getenv("PUBSUB_PROJECT_ID"),
		},
		watermillLogger,
	)
	if err != nil {
		appDependencies.Logger.Error("failed to connect to Google Pub/Sub subscriber",
			slog.String("error", err.Error()),
		)
		panic(err)
	}

	ssepubsub.Routes(router, subscriber, appDependencies)

	return router
}
