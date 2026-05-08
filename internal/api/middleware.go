package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (handler *Handler) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusWriter := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(statusWriter, r)
		handler.deps.Logger.Info("http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", statusWriter.status),
			slog.String("duration", time.Since(start).String()),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (writer *loggingResponseWriter) WriteHeader(status int) {
	writer.status = status
	writer.ResponseWriter.WriteHeader(status)
}
