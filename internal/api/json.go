package api

import (
	"encoding/json"
	"net/http"

	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
)

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeImportError(w http.ResponseWriter, status int, err mailbox.ImportError) {
	writeJSON(w, status, ErrorResponse{
		Error:   err.What,
		Kind:    err.Kind,
		What:    err.What,
		Why:     err.Why,
		NowWhat: err.NowWhat,
	})
}
