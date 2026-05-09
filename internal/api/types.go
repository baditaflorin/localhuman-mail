package api

import (
	"log/slog"

	"github.com/baditaflorin/localhuman-mail/internal/ai"
	"github.com/baditaflorin/localhuman-mail/internal/config"
	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
	"github.com/baditaflorin/localhuman-mail/internal/metrics"
	"github.com/baditaflorin/localhuman-mail/internal/search"
	"github.com/baditaflorin/localhuman-mail/internal/security"
)

type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

type Dependencies struct {
	Config    config.Config
	Store     *search.Store
	Assistant *ai.Service
	Metrics   *metrics.Recorder
	Logger    *slog.Logger
	Version   VersionInfo
}

type HealthResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Kind    string `json:"kind,omitempty"`
	What    string `json:"what,omitempty"`
	Why     string `json:"why,omitempty"`
	NowWhat string `json:"nowWhat,omitempty"`
}

type capabilitiesResponse struct {
	Capabilities []security.Capability `json:"capabilities"`
}

type messageListResponse struct {
	Messages []mailbox.Message `json:"messages"`
}

type searchResponse struct {
	Query     string            `json:"query"`
	ElapsedMS int64             `json:"elapsedMs"`
	Messages  []mailbox.Message `json:"messages"`
}

type assistReplyRequest struct {
	MessageID    string `json:"messageId" validate:"required"`
	Tone         string `json:"tone" validate:"required,oneof=concise warm decisive"`
	Instructions string `json:"instructions" validate:"max=1200"`
}
