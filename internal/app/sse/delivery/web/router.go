package web

import (
	"github.com/Christian-007/fit-forge-notification/internal/pkg/appdependency"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/appservices"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes(appDependencies appdependency.AppDependency) *chi.Mux {
	r := chi.NewRouter()

	sseHandler := NewSseHandler(SseHandlerOptions{
		appDependencies.Logger,
		appDependencies.InMemoryMessageBroker,
	})

	authService := appservices.NewAuthServiceImpl(appservices.AuthServiceOptions{
		Cache: appDependencies.RedisClient,
	})
	strictSessionMiddleware := middlewares.StrictSession(authService, appDependencies.SecretManagerClient)

	r.Group(func(r chi.Router) {
		r.Use(strictSessionMiddleware)
		r.Get("/rewards", sseHandler.GetRewards)
	})

	return r
}
