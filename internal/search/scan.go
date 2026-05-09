package search

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
)

var ErrNotFound = errors.New("message not found")

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMessages(rows *sql.Rows) ([]mailbox.Message, error) {
	messages := make([]mailbox.Message, 0)
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}
	return messages, nil
}

func scanMessage(row rowScanner) (mailbox.Message, error) {
	var message mailbox.Message
	var toJSON, tagsJSON, dateValue string
	var confidenceJSON, fieldConfidenceJSON, warningsJSON, attachmentsJSON, calendarJSON, provenanceJSON string
	if err := row.Scan(
		&message.ID,
		&message.SourceID,
		&message.Subject,
		&message.From,
		&toJSON,
		&dateValue,
		&message.Snippet,
		&message.Body,
		&message.PrimaryBody,
		&message.Shape,
		&tagsJSON,
		&confidenceJSON,
		&fieldConfidenceJSON,
		&warningsJSON,
		&attachmentsJSON,
		&calendarJSON,
		&provenanceJSON,
	); err != nil {
		return mailbox.Message{}, err
	}
	if err := json.Unmarshal([]byte(toJSON), &message.To); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal recipients: %w", err)
	}
	if err := json.Unmarshal([]byte(tagsJSON), &message.Tags); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal tags: %w", err)
	}
	if err := json.Unmarshal([]byte(confidenceJSON), &message.Confidence); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal confidence: %w", err)
	}
	if err := json.Unmarshal([]byte(fieldConfidenceJSON), &message.FieldConfidence); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal field confidence: %w", err)
	}
	if err := json.Unmarshal([]byte(warningsJSON), &message.Warnings); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal warnings: %w", err)
	}
	if err := json.Unmarshal([]byte(attachmentsJSON), &message.Attachments); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal attachments: %w", err)
	}
	if strings.TrimSpace(calendarJSON) != "" && strings.TrimSpace(calendarJSON) != "null" {
		var calendar mailbox.CalendarEvent
		if err := json.Unmarshal([]byte(calendarJSON), &calendar); err != nil {
			return mailbox.Message{}, fmt.Errorf("unmarshal calendar: %w", err)
		}
		message.Calendar = &calendar
	}
	if err := json.Unmarshal([]byte(provenanceJSON), &message.Provenance); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal provenance: %w", err)
	}
	date, err := time.Parse(time.RFC3339Nano, dateValue)
	if err != nil {
		return mailbox.Message{}, fmt.Errorf("parse message date: %w", err)
	}
	message.Date = date
	if message.PrimaryBody == "" {
		message.PrimaryBody = message.Body
	}
	if message.Shape == "" {
		message.Shape = mailbox.UnknownShape
	}
	if message.Confidence.Score == 0 {
		message.Confidence = mailbox.NewConfidence(0.55, "legacy message without stored confidence")
	}
	if message.FieldConfidence == nil {
		message.FieldConfidence = map[string]mailbox.Confidence{}
	}
	if message.Warnings == nil {
		message.Warnings = []mailbox.ImportWarning{}
	}
	if message.Attachments == nil {
		message.Attachments = []mailbox.Attachment{}
	}
	if message.SourceID == "" {
		message.SourceID = message.ID
	}
	return message, nil
}

func saneLimit(limit int) int {
	if limit < 1 {
		return 50
	}
	if limit > 200 {
		return 200
	}
	return limit
}
