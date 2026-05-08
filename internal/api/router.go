package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	deps     Dependencies
	validate *validator.Validate
}

func NewRouter(deps Dependencies) http.Handler {
	handler := &Handler{
		deps:     deps,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   deps.Config.AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(handler.logRequests)
	router.Use(deps.Metrics.Middleware)
	router.Use(middleware.Recoverer)

	router.Get("/healthz", handler.healthz)
	router.Get("/readyz", handler.readyz)
	router.Handle("/metrics", deps.Metrics.Handler())

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/version", handler.version)
		r.Get("/capabilities", handler.capabilities)
		r.Get("/messages", handler.listMessages)
		r.Get("/messages/{id}", handler.getMessage)
		r.Get("/search", handler.searchMessages)
		r.Post("/import/demo", handler.importDemo)
		r.Post("/import/eml", handler.importEML)
		r.Post("/assist/reply", handler.assistReply)
	})

	return router
}
