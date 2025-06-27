package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Christian-007/fit-forge-notification/internal/pkg/appdependency"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/cache"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/messagebroker"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/security"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

const Addr = ":4001"

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a connection to a Message Broker
	watermillLogger := watermill.NewStdLogger(false, false)

	// Open Redis Connection
	redisClient, err := cache.NewRedisCache(&redis.Options{
		Addr:     os.Getenv("REDIS_DSN"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err != nil {
		logger.Error("Failed to connect to Redis",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Create GCP Secret Manager Client
	var secretManagerClient security.SecretManageProvider
	if os.Getenv("ENV") == "production" {
		secretManagerClient, err = security.NewGCPSecretManagerClient(context.Background())
		if err != nil {
			logger.Error("failed to create GCP Secret Manager Client",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	} else {
		secretManagerClient, err = security.NewLocalSecretManager()
		if err != nil {
			logger.Error("failed to create Local Secret Manager",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}
	defer secretManagerClient.Close()

	// Instantiate the all application dependencies
	appDependencies := appdependency.NewAppDependency(appdependency.AppDependencyOptions{
		Logger:                logger,
		InMemoryMessageBroker: messagebroker.NewInMemoryMessageBroker(logger),
		RedisClient:           redisClient,
		SecretManagerClient:   secretManagerClient,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errgrp, ctx := errgroup.WithContext(ctx)

	// Instantiate the PubSub router
	watermillRouter := NewWatermillRouter(watermillLogger, appDependencies)
	errgrp.Go(func() error {
		logger.Info("starting PubSub router...")
		err := watermillRouter.Run(ctx) // Starting the PubSub router in a Goroutine
		if err != nil {
			logger.Error("failed to start PubSub router",
				slog.String("error", err.Error()),
			)

			return err
		}
		return nil
	})

	// HTTP Server configurations (Non TLS)
	server := &http.Server{
		Addr:        Addr,
		Handler:     Routes(appDependencies),
		ErrorLog:    slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
	}

	errgrp.Go(func() error {
		// We don't want to start the HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-watermillRouter.Running()

		logger.Info("starting server", "addr", Addr)

		err := server.ListenAndServe()
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(ctx)
	})

	err = errgrp.Wait()
	if err != nil {
		panic(err)
	}
}
