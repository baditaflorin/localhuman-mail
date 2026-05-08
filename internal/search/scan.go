package search

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	if err := row.Scan(
		&message.ID,
		&message.Subject,
		&message.From,
		&toJSON,
		&dateValue,
		&message.Snippet,
		&message.Body,
		&tagsJSON,
	); err != nil {
		return mailbox.Message{}, err
	}
	if err := json.Unmarshal([]byte(toJSON), &message.To); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal recipients: %w", err)
	}
	if err := json.Unmarshal([]byte(tagsJSON), &message.Tags); err != nil {
		return mailbox.Message{}, fmt.Errorf("unmarshal tags: %w", err)
	}
	date, err := time.Parse(time.RFC3339Nano, dateValue)
	if err != nil {
		return mailbox.Message{}, fmt.Errorf("parse message date: %w", err)
	}
	message.Date = date
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
