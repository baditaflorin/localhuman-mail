package mailbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type fixtureExpectation struct {
	Source                  string   `json:"source"`
	MustImport              bool     `json:"mustImport"`
	ExpectedShape           string   `json:"expectedShape"`
	MustContain             []string `json:"mustContain"`
	MinConfidence           float64  `json:"minConfidence"`
	ExpectedAttachmentCount int      `json:"expectedAttachmentCount"`
	ExpectedWarningsContain []string `json:"expectedWarningsContain"`
	Deterministic           bool     `json:"deterministic"`
}

func TestRealDataFixtures(t *testing.T) {
	fixturesDir := filepath.Join("..", "..", "test", "fixtures", "realdata")
	expectedFiles, err := filepath.Glob(filepath.Join(fixturesDir, "*.expected.json"))
	require.NoError(t, err)
	require.NotEmpty(t, expectedFiles)

	for _, expectedPath := range expectedFiles {
		name := strings.TrimSuffix(filepath.Base(expectedPath), ".expected.json")
		t.Run(name, func(t *testing.T) {
			expectation := readFixtureExpectation(t, expectedPath)
			raw, err := os.ReadFile(filepath.Join(fixturesDir, name+".eml"))
			require.NoError(t, err)

			message, parseErr := ParseEML(bytes.NewReader(raw))
			if !expectation.MustImport {
				require.Error(t, parseErr)
				return
			}

			require.NoError(t, parseErr, expectation.Source)
			require.Equal(t, expectation.ExpectedShape, message.Shape)
			require.GreaterOrEqual(t, message.Confidence.Score, expectation.MinConfidence)
			require.Len(t, message.Attachments, expectation.ExpectedAttachmentCount)

			searchable := strings.ToLower(strings.Join([]string{
				message.Subject,
				message.From,
				message.Body,
				message.PrimaryBody,
				calendarText(message.Calendar),
				warningsText(message.Warnings),
			}, " "))
			for _, required := range expectation.MustContain {
				require.Contains(t, searchable, strings.ToLower(required), "fixture must contain expected domain text")
			}
			for _, requiredWarning := range expectation.ExpectedWarningsContain {
				require.Contains(t, warningsText(message.Warnings), strings.ToLower(requiredWarning), "fixture must surface expected warning")
			}
			if expectation.Deterministic {
				second, err := ParseEML(bytes.NewReader(raw))
				require.NoError(t, err)
				firstJSON, err := json.Marshal(message)
				require.NoError(t, err)
				secondJSON, err := json.Marshal(second)
				require.NoError(t, err)
				require.JSONEq(t, string(firstJSON), string(secondJSON))
			}
		})
	}
}

func TestOversizeEmailIsActionable(t *testing.T) {
	raw := bytes.Repeat([]byte("x"), MaxEMLBytes+1)

	_, err := ParseEML(bytes.NewReader(raw))

	var importErr ImportError
	require.ErrorAs(t, err, &importErr)
	require.Equal(t, "Email is too large", importErr.What)
	require.NotEmpty(t, importErr.NowWhat)
}

func TestUnreadableEmailIsActionable(t *testing.T) {
	_, err := ParseEML(strings.NewReader("this is not an email"))

	var importErr ImportError
	require.True(t, errors.As(err, &importErr))
	require.Equal(t, "Email could not be opened", importErr.What)
	require.NotEmpty(t, importErr.Why)
	require.NotEmpty(t, importErr.NowWhat)
}

func readFixtureExpectation(t *testing.T, path string) fixtureExpectation {
	t.Helper()

	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var expectation fixtureExpectation
	require.NoError(t, json.Unmarshal(raw, &expectation))
	return expectation
}

func warningsText(warnings []ImportWarning) string {
	parts := make([]string, 0, len(warnings)*3)
	for _, warning := range warnings {
		parts = append(parts, warning.Field, warning.Message, warning.NextStep)
	}
	return strings.ToLower(strings.Join(parts, " "))
}
