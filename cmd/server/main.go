package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/ai"
	"github.com/baditaflorin/localhuman-mail/internal/api"
	"github.com/baditaflorin/localhuman-mail/internal/config"
	"github.com/baditaflorin/localhuman-mail/internal/metrics"
	"github.com/baditaflorin/localhuman-mail/internal/search"
)

var (
	version   = "0.3.1"
	commit    = "dev"
	buildTime = "unknown"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "check the local HTTP health endpoint")
	flag.Parse()
	if *healthcheck {
		runHealthcheck()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error("load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel.Level(),
	}))
	slog.SetDefault(logger)

	store, err := search.Open(ctx, cfg.DataDir)
	if err != nil {
		logger.Error("open store", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	recorder := metrics.NewRecorder()
	assistant := ai.NewService(cfg)
	router := api.NewRouter(api.Dependencies{
		Config:    cfg,
		Store:     store,
		Assistant: assistant,
		Metrics:   recorder,
		Logger:    logger,
		Version: api.VersionInfo{
			Version:   version,
			Commit:    commit,
			BuildTime: buildTime,
		},
	})

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server listening", "addr", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown", "error", err)
	}
}

func runHealthcheck() {
	client := http.Client{Timeout: 2 * time.Second}
	response, err := client.Get("http://127.0.0.1:8080/healthz")
	if err != nil {
		os.Exit(1)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
