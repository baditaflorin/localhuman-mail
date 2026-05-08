package search

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	db, err := sql.Open("sqlite", filepath.Join(dataDir, "localhuman-mail.db"))
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	store := &Store{db: db}
	if err := store.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) Ready(ctx context.Context) error {
	return store.db.PingContext(ctx)
}

func (store *Store) AddMessages(ctx context.Context, messages []mailbox.Message) (mailbox.ImportResult, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return mailbox.ImportResult{}, fmt.Errorf("begin import: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	result := mailbox.ImportResult{Total: len(messages)}
	for _, message := range messages {
		toJSON, tagsJSON, err := marshalLists(message)
		if err != nil {
			return mailbox.ImportResult{}, err
		}
		sqlResult, err := tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO messages
				(id, subject, sender, recipients_json, date, snippet, body, tags_json, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, message.ID, message.Subject, message.From, toJSON, message.Date.UTC().Format(time.RFC3339Nano), message.Snippet, message.Body, tagsJSON, time.Now().UTC().Format(time.RFC3339Nano))
		if err != nil {
			return mailbox.ImportResult{}, fmt.Errorf("insert message %s: %w", message.ID, err)
		}
		affected, err := sqlResult.RowsAffected()
		if err != nil {
			return mailbox.ImportResult{}, fmt.Errorf("message insert rows affected: %w", err)
		}
		if affected == 0 {
			result.Skipped++
			continue
		}
		result.Imported++
	}
	if err = tx.Commit(); err != nil {
		return mailbox.ImportResult{}, fmt.Errorf("commit import: %w", err)
	}
	return result, nil
}

func (store *Store) ListMessages(ctx context.Context, limit int) ([]mailbox.Message, error) {
	rows, err := store.db.QueryContext(ctx, `
		SELECT id, subject, sender, recipients_json, date, snippet, body, tags_json
		FROM messages
		ORDER BY date DESC
		LIMIT ?
	`, saneLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()
	return scanMessages(rows)
}

func (store *Store) GetMessage(ctx context.Context, id string) (mailbox.Message, error) {
	row := store.db.QueryRowContext(ctx, `
		SELECT id, subject, sender, recipients_json, date, snippet, body, tags_json
		FROM messages
		WHERE id = ?
	`, id)
	message, err := scanMessage(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mailbox.Message{}, ErrNotFound
		}
		return mailbox.Message{}, fmt.Errorf("get message: %w", err)
	}
	return message, nil
}

func (store *Store) SearchMessages(ctx context.Context, query string, limit int) ([]mailbox.Message, error) {
	like := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
	rows, err := store.db.QueryContext(ctx, `
		SELECT id, subject, sender, recipients_json, date, snippet, body, tags_json
		FROM messages
		WHERE lower(subject) LIKE ?
		   OR lower(sender) LIKE ?
		   OR lower(snippet) LIKE ?
		   OR lower(body) LIKE ?
		   OR lower(tags_json) LIKE ?
		ORDER BY date DESC
		LIMIT ?
	`, like, like, like, like, like, saneLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}
	defer rows.Close()
	return scanMessages(rows)
}

func (store *Store) migrate(ctx context.Context) error {
	_, err := store.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			subject TEXT NOT NULL,
			sender TEXT NOT NULL,
			recipients_json TEXT NOT NULL,
			date TEXT NOT NULL,
			snippet TEXT NOT NULL,
			body TEXT NOT NULL,
			tags_json TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_messages_date ON messages(date);
		CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender);
	`)
	if err != nil {
		return fmt.Errorf("migrate sqlite: %w", err)
	}
	return nil
}

func marshalLists(message mailbox.Message) (string, string, error) {
	toJSON, err := json.Marshal(message.To)
	if err != nil {
		return "", "", fmt.Errorf("marshal recipients: %w", err)
	}
	tagsJSON, err := json.Marshal(message.Tags)
	if err != nil {
		return "", "", fmt.Errorf("marshal tags: %w", err)
	}
	return string(toJSON), string(tagsJSON), nil
}
