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

	// Register the pure-Go SQLite driver for database/sql.
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
		values, err := marshalMessageJSON(message)
		if err != nil {
			return mailbox.ImportResult{}, err
		}
		sqlResult, err := tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO messages
				(id, source_id, subject, sender, recipients_json, date, snippet, body, primary_body, shape, tags_json, confidence_json, field_confidence_json, warnings_json, attachments_json, calendar_json, provenance_json, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			message.ID,
			message.SourceID,
			message.Subject,
			message.From,
			values.recipients,
			message.Date.UTC().Format(time.RFC3339Nano),
			message.Snippet,
			message.Body,
			message.PrimaryBody,
			message.Shape,
			values.tags,
			values.confidence,
			values.fieldConfidence,
			values.warnings,
			values.attachments,
			values.calendar,
			values.provenance,
			time.Now().UTC().Format(time.RFC3339Nano),
		)
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
		SELECT id, source_id, subject, sender, recipients_json, date, snippet, body, primary_body, shape, tags_json, confidence_json, field_confidence_json, warnings_json, attachments_json, calendar_json, provenance_json
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
		SELECT id, source_id, subject, sender, recipients_json, date, snippet, body, primary_body, shape, tags_json, confidence_json, field_confidence_json, warnings_json, attachments_json, calendar_json, provenance_json
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
		SELECT id, source_id, subject, sender, recipients_json, date, snippet, body, primary_body, shape, tags_json, confidence_json, field_confidence_json, warnings_json, attachments_json, calendar_json, provenance_json
		FROM messages
		WHERE lower(subject) LIKE ?
		   OR lower(sender) LIKE ?
		   OR lower(snippet) LIKE ?
		   OR lower(body) LIKE ?
		   OR lower(primary_body) LIKE ?
		   OR lower(shape) LIKE ?
		   OR lower(tags_json) LIKE ?
		   OR lower(warnings_json) LIKE ?
		   OR lower(attachments_json) LIKE ?
		ORDER BY date DESC
		LIMIT ?
	`, like, like, like, like, like, like, like, like, like, saneLimit(limit))
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
			source_id TEXT NOT NULL DEFAULT '',
			subject TEXT NOT NULL,
			sender TEXT NOT NULL,
			recipients_json TEXT NOT NULL,
			date TEXT NOT NULL,
			snippet TEXT NOT NULL,
			body TEXT NOT NULL,
			primary_body TEXT NOT NULL DEFAULT '',
			shape TEXT NOT NULL DEFAULT '',
			tags_json TEXT NOT NULL,
			confidence_json TEXT NOT NULL DEFAULT '{}',
			field_confidence_json TEXT NOT NULL DEFAULT '{}',
			warnings_json TEXT NOT NULL DEFAULT '[]',
			attachments_json TEXT NOT NULL DEFAULT '[]',
			calendar_json TEXT NOT NULL DEFAULT 'null',
			provenance_json TEXT NOT NULL DEFAULT '{}',
			created_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_messages_date ON messages(date);
		CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender);
	`)
	if err != nil {
		return fmt.Errorf("migrate sqlite: %w", err)
	}
	columns := map[string]string{
		"source_id":             "ALTER TABLE messages ADD COLUMN source_id TEXT NOT NULL DEFAULT ''",
		"primary_body":          "ALTER TABLE messages ADD COLUMN primary_body TEXT NOT NULL DEFAULT ''",
		"shape":                 "ALTER TABLE messages ADD COLUMN shape TEXT NOT NULL DEFAULT ''",
		"confidence_json":       "ALTER TABLE messages ADD COLUMN confidence_json TEXT NOT NULL DEFAULT '{}'",
		"field_confidence_json": "ALTER TABLE messages ADD COLUMN field_confidence_json TEXT NOT NULL DEFAULT '{}'",
		"warnings_json":         "ALTER TABLE messages ADD COLUMN warnings_json TEXT NOT NULL DEFAULT '[]'",
		"attachments_json":      "ALTER TABLE messages ADD COLUMN attachments_json TEXT NOT NULL DEFAULT '[]'",
		"calendar_json":         "ALTER TABLE messages ADD COLUMN calendar_json TEXT NOT NULL DEFAULT 'null'",
		"provenance_json":       "ALTER TABLE messages ADD COLUMN provenance_json TEXT NOT NULL DEFAULT '{}'",
	}
	for name, ddl := range columns {
		if err := store.ensureColumn(ctx, name, ddl); err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) ensureColumn(ctx context.Context, name string, ddl string) error {
	rows, err := store.db.QueryContext(ctx, "PRAGMA table_info(messages)")
	if err != nil {
		return fmt.Errorf("inspect messages schema: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var columnName, columnType string
		var notNull, primaryKey int
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &columnName, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return fmt.Errorf("scan messages schema: %w", err)
		}
		if columnName == name {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate messages schema: %w", err)
	}
	if _, err := store.db.ExecContext(ctx, ddl); err != nil {
		return fmt.Errorf("migrate messages column %s: %w", name, err)
	}
	return nil
}

type messageJSON struct {
	recipients      string
	tags            string
	confidence      string
	fieldConfidence string
	warnings        string
	attachments     string
	calendar        string
	provenance      string
}

func marshalMessageJSON(message mailbox.Message) (messageJSON, error) {
	result := messageJSON{}
	toJSON, err := json.Marshal(message.To)
	if err != nil {
		return result, fmt.Errorf("marshal recipients: %w", err)
	}
	result.recipients = string(toJSON)
	tagsJSON, err := json.Marshal(message.Tags)
	if err != nil {
		return result, fmt.Errorf("marshal tags: %w", err)
	}
	result.tags = string(tagsJSON)
	confidenceJSON, err := json.Marshal(message.Confidence)
	if err != nil {
		return result, fmt.Errorf("marshal confidence: %w", err)
	}
	result.confidence = string(confidenceJSON)
	fieldConfidenceJSON, err := json.Marshal(message.FieldConfidence)
	if err != nil {
		return result, fmt.Errorf("marshal field confidence: %w", err)
	}
	result.fieldConfidence = string(fieldConfidenceJSON)
	warningsJSON, err := json.Marshal(message.Warnings)
	if err != nil {
		return result, fmt.Errorf("marshal warnings: %w", err)
	}
	result.warnings = string(warningsJSON)
	attachmentsJSON, err := json.Marshal(message.Attachments)
	if err != nil {
		return result, fmt.Errorf("marshal attachments: %w", err)
	}
	result.attachments = string(attachmentsJSON)
	calendarJSON, err := json.Marshal(message.Calendar)
	if err != nil {
		return result, fmt.Errorf("marshal calendar: %w", err)
	}
	result.calendar = string(calendarJSON)
	provenanceJSON, err := json.Marshal(message.Provenance)
	if err != nil {
		return result, fmt.Errorf("marshal provenance: %w", err)
	}
	result.provenance = string(provenanceJSON)
	return result, nil
}
