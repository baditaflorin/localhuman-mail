package api

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/baditaflorin/localhuman-mail/internal/ai"
	"github.com/baditaflorin/localhuman-mail/internal/config"
	"github.com/baditaflorin/localhuman-mail/internal/metrics"
	"github.com/baditaflorin/localhuman-mail/internal/search"
	"github.com/stretchr/testify/require"
)

func TestRouterHealthAndDemoImport(t *testing.T) {
	cfg := config.Config{
		Addr:           ":0",
		DataDir:        t.TempDir(),
		AllowedOrigins: []string{"http://localhost:5173"},
		LogLevel:       "error",
		OllamaURL:      "",
	}
	store, err := search.Open(context.Background(), cfg.DataDir)
	require.NoError(t, err)
	defer store.Close()

	router := NewRouter(Dependencies{
		Config:    cfg,
		Store:     store,
		Assistant: ai.NewService(cfg),
		Metrics:   metrics.NewRecorder(),
		Logger:    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Version:   VersionInfo{Version: "test", Commit: "test", BuildTime: "test"},
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/import/demo", nil)
	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), `"imported"`)
}
