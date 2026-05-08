package mailbox

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEML(t *testing.T) {
	raw := strings.NewReader("From: Maya <maya@example.com>\r\nTo: You <you@example.com>\r\nSubject: Test mail\r\nDate: Fri, 08 May 2026 10:00:00 +0000\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nHello from localhuman.")

	message, err := ParseEML(raw)

	require.NoError(t, err)
	require.Equal(t, "Test mail", message.Subject)
	require.Contains(t, message.From, "maya@example.com")
	require.Contains(t, message.Body, "Hello from localhuman")
	require.NotEmpty(t, message.ID)
}
