package mailbox

import (
	"crypto/rand"
	"encoding/base64"
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

// TestParseEMLReportsTrueAttachmentSizeBeyondIndexCap guards against a
// regression where attachment SizeBytes was silently capped at
// MaxIndexedBytes (2MB) because countPart reused the same limit reader used
// for indexed body text. Attachment bytes are never retained (only
// metadata), so the reported size must reflect the real attachment size even
// when it exceeds the indexed-text budget.
func TestParseEMLReportsTrueAttachmentSizeBeyondIndexCap(t *testing.T) {
	rawAttachment := make([]byte, MaxIndexedBytes+64*1024) // exceed the 2MB indexed-text cap
	_, err := rand.Read(rawAttachment)
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(rawAttachment)

	var body strings.Builder
	body.WriteString("From: sender@example.com\r\n")
	body.WriteString("To: receiver@example.com\r\n")
	body.WriteString("Subject: big attachment\r\n")
	body.WriteString("Date: Fri, 08 May 2026 10:00:00 +0000\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: multipart/mixed; boundary=\"BOUNDARY123\"\r\n")
	body.WriteString("\r\n--BOUNDARY123\r\n")
	body.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	body.WriteString("See attached file.\r\n")
	body.WriteString("\r\n--BOUNDARY123\r\n")
	body.WriteString("Content-Type: application/octet-stream; name=\"big.bin\"\r\n")
	body.WriteString("Content-Transfer-Encoding: base64\r\n")
	body.WriteString("Content-Disposition: attachment; filename=\"big.bin\"\r\n\r\n")
	for encoded != "" {
		chunkLen := 76
		if len(encoded) < chunkLen {
			chunkLen = len(encoded)
		}
		body.WriteString(encoded[:chunkLen])
		body.WriteString("\r\n")
		encoded = encoded[chunkLen:]
	}
	body.WriteString("--BOUNDARY123--\r\n")

	message, err := ParseEML(strings.NewReader(body.String()))
	require.NoError(t, err)
	require.Len(t, message.Attachments, 1)
	require.Equal(t, int64(len(rawAttachment)), message.Attachments[0].SizeBytes)
}
