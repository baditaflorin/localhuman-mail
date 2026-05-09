package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
	"github.com/baditaflorin/localhuman-mail/internal/search"
	"github.com/baditaflorin/localhuman-mail/internal/security"
	"github.com/go-chi/chi/v5"
)

func (handler *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (handler *Handler) readyz(w http.ResponseWriter, r *http.Request) {
	if err := handler.deps.Store.Ready(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "store is not ready")
		return
	}
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (handler *Handler) version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, handler.deps.Version)
}

func (handler *Handler) capabilities(w http.ResponseWriter, r *http.Request) {
	capabilities := security.DetectCapabilities(r.Context(), handler.deps.Config)
	writeJSON(w, http.StatusOK, capabilitiesResponse{Capabilities: capabilities})
}

func (handler *Handler) listMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := handler.deps.Store.ListMessages(r.Context(), limitParam(r, 50))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list messages")
		return
	}
	writeJSON(w, http.StatusOK, messageListResponse{Messages: messages})
}

func (handler *Handler) getMessage(w http.ResponseWriter, r *http.Request) {
	message, err := handler.deps.Store.GetMessage(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, search.ErrNotFound) {
			writeError(w, http.StatusNotFound, "message not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not load message")
		return
	}
	writeJSON(w, http.StatusOK, message)
}

func (handler *Handler) searchMessages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "q is required")
		return
	}
	handler.deps.Metrics.AddSearchRequest()
	messages, err := handler.deps.Store.SearchMessages(r.Context(), query, limitParam(r, 25))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not search messages")
		return
	}
	writeJSON(w, http.StatusOK, searchResponse{
		Query:     query,
		ElapsedMS: time.Since(start).Milliseconds(),
		Messages:  messages,
	})
}

func (handler *Handler) importDemo(w http.ResponseWriter, r *http.Request) {
	result, err := handler.deps.Store.AddMessages(r.Context(), mailbox.DemoMessages(time.Now().UTC()))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not import demo messages")
		return
	}
	handler.deps.Metrics.AddImported(result.Imported)
	writeJSON(w, http.StatusOK, result)
}

func (handler *Handler) importEML(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, mailbox.MaxEMLBytes+(1<<20))
	// #nosec G120 -- the request body is capped above before multipart parsing.
	if err := r.ParseMultipartForm(mailbox.MaxEMLBytes + (1 << 20)); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart upload")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	message, err := mailbox.ParseEML(file)
	if err != nil {
		var importErr mailbox.ImportError
		if errors.As(err, &importErr) {
			writeImportError(w, http.StatusBadRequest, importErr)
			return
		}
		writeError(w, http.StatusBadRequest, "could not parse EML")
		return
	}
	result, err := handler.deps.Store.AddMessages(r.Context(), []mailbox.Message{message})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not store message")
		return
	}
	handler.deps.Metrics.AddImported(result.Imported)
	writeJSON(w, http.StatusOK, result)
}

func (handler *Handler) assistReply(w http.ResponseWriter, r *http.Request) {
	var request assistReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := handler.validate.Struct(request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid reply assist request")
		return
	}
	message, err := handler.deps.Store.GetMessage(r.Context(), request.MessageID)
	if err != nil {
		if errors.Is(err, search.ErrNotFound) {
			writeError(w, http.StatusNotFound, "message not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not load message")
		return
	}
	handler.deps.Metrics.AddAssistRequest()
	writeJSON(w, http.StatusOK, handler.deps.Assistant.AssistReply(r.Context(), message, request.Tone, request.Instructions))
}

func limitParam(r *http.Request, fallback int) int {
	value := r.URL.Query().Get("limit")
	if value == "" {
		return fallback
	}
	limit, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	if limit < 1 {
		return fallback
	}
	if limit > 200 {
		return 200
	}
	return limit
}
