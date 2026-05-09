package mailbox

import (
	"math"
	"time"
)

const (
	ParserVersion    = "phase2-parser-v1"
	SchemaVersion    = "message.v2"
	MaxEMLBytes      = 25 * 1024 * 1024
	MaxIndexedBytes  = 2 * 1024 * 1024
	UnknownShape     = "unknown"
	LowConfidence    = "low"
	MediumConfidence = "medium"
	HighConfidence   = "high"
)

type Message struct {
	ID              string                `json:"id"`
	SourceID        string                `json:"sourceId"`
	Subject         string                `json:"subject"`
	From            string                `json:"from"`
	To              []string              `json:"to"`
	Date            time.Time             `json:"date"`
	Snippet         string                `json:"snippet"`
	Body            string                `json:"body"`
	PrimaryBody     string                `json:"primaryBody"`
	Shape           string                `json:"shape"`
	Tags            []string              `json:"tags"`
	Confidence      Confidence            `json:"confidence"`
	FieldConfidence map[string]Confidence `json:"fieldConfidence"`
	Warnings        []ImportWarning       `json:"warnings"`
	Attachments     []Attachment          `json:"attachments"`
	Calendar        *CalendarEvent        `json:"calendar,omitempty"`
	Provenance      Provenance            `json:"provenance"`
}

type ImportResult struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
	Total    int `json:"total"`
}

type Confidence struct {
	Score   float64  `json:"score"`
	Label   string   `json:"label"`
	Reasons []string `json:"reasons"`
}

type ImportWarning struct {
	Severity string `json:"severity"`
	Field    string `json:"field"`
	Message  string `json:"message"`
	NextStep string `json:"nextStep"`
}

type Attachment struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	SizeBytes   int64  `json:"sizeBytes"`
}

type CalendarEvent struct {
	Summary  string `json:"summary"`
	Location string `json:"location"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

type Provenance struct {
	SourceID      string `json:"sourceId"`
	SourceSHA256  string `json:"sourceSha256"`
	ParserVersion string `json:"parserVersion"`
	SchemaVersion string `json:"schemaVersion"`
	SizeBytes     int64  `json:"sizeBytes"`
}

type ImportError struct {
	Kind    string `json:"kind"`
	What    string `json:"what"`
	Why     string `json:"why"`
	NowWhat string `json:"nowWhat"`
}

func (err ImportError) Error() string {
	return err.What + ": " + err.Why + " " + err.NowWhat
}

func NewConfidence(score float64, reasons ...string) Confidence {
	score = math.Round(score*100) / 100
	label := HighConfidence
	switch {
	case score < 0.5:
		label = LowConfidence
	case score < 0.8:
		label = MediumConfidence
	}
	return Confidence{Score: score, Label: label, Reasons: reasons}
}
