package main

import (
	"net/http"
	"os"

	sseweb "github.com/Christian-007/fit-forge-notification/internal/app/sse/delivery/web"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/appdependency"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/middlewares"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(appDependencies appdependency.AppDependency) *chi.Mux {
	r := chi.NewRouter()

	logRequest := middlewares.NewLogRequest(appDependencies.Logger)

	r.Use(logRequest)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONTEND_URL")},
		AllowCredentials: true,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.SendResponse(w, http.StatusOK, utils.ErrorResponse{Message: "Ok"})
	})

	r.Mount("/sse", sseweb.Routes(appDependencies))

	return r
}
