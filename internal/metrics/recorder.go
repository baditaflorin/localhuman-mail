package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Recorder struct {
	registry       *prometheus.Registry
	httpRequests   *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
	imported       prometheus.Counter
	searchRequests prometheus.Counter
	assistRequests prometheus.Counter
}

func NewRecorder() *Recorder {
	registry := prometheus.NewRegistry()
	recorder := &Recorder{
		registry: registry,
		httpRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "localhuman_http_requests_total",
			Help: "HTTP requests by method, route, and status.",
		}, []string{"method", "route", "status"}),
		httpDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "localhuman_http_request_duration_seconds",
			Help:    "HTTP request duration by method and route.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "route"}),
		imported: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "localhuman_messages_imported_total",
			Help: "Imported messages.",
		}),
		searchRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "localhuman_search_requests_total",
			Help: "Search requests.",
		}),
		assistRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "localhuman_assist_requests_total",
			Help: "AI assist requests.",
		}),
	}
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(recorder.httpRequests, recorder.httpDuration, recorder.imported, recorder.searchRequests, recorder.assistRequests)
	return recorder
}

func (recorder *Recorder) Handler() http.Handler {
	return promhttp.HandlerFor(recorder.registry, promhttp.HandlerOpts{})
}

func (recorder *Recorder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusWriter := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(statusWriter, r)
		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			route = r.URL.Path
		}
		status := strconv.Itoa(statusWriter.status)
		recorder.httpRequests.WithLabelValues(r.Method, route, status).Inc()
		recorder.httpDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

func (recorder *Recorder) AddImported(count int) {
	if count > 0 {
		recorder.imported.Add(float64(count))
	}
}

func (recorder *Recorder) AddSearchRequest() {
	recorder.searchRequests.Inc()
}

func (recorder *Recorder) AddAssistRequest() {
	recorder.assistRequests.Inc()
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (writer *statusResponseWriter) WriteHeader(status int) {
	writer.status = status
	writer.ResponseWriter.WriteHeader(status)
}
