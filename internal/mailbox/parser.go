package mailbox

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"strings"
	"time"

	"github.com/emersion/go-message/mail"
)

func ParseEML(reader io.Reader) (Message, error) {
	mailReader, err := mail.CreateReader(reader)
	if err != nil {
		return Message{}, err
	}

	subject, _ := mailReader.Header.Subject()
	date, err := mailReader.Header.Date()
	if err != nil {
		date = time.Now().UTC()
	}

	from := addressLine(mailReader.Header, "From")
	to := addressList(mailReader.Header, "To")
	body, err := readPlainBody(mailReader)
	if err != nil {
		return Message{}, err
	}
	if strings.TrimSpace(body) == "" {
		return Message{}, errors.New("message has no readable text body")
	}

	message := Message{
		ID:      messageID(subject, from, date, body),
		Subject: fallback(subject, "(no subject)"),
		From:    fallback(from, "unknown"),
		To:      to,
		Date:    date.UTC(),
		Snippet: snippet(body),
		Body:    strings.TrimSpace(body),
		Tags:    []string{"imported"},
	}
	return message, nil
}

func readPlainBody(mailReader *mail.Reader) (string, error) {
	var parts []string
	for {
		part, err := mailReader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		header, ok := part.Header.(*mail.InlineHeader)
		if !ok {
			continue
		}
		contentType, _, err := header.ContentType()
		if err != nil {
			contentType = "text/plain"
		}
		mediaType, _, _ := mime.ParseMediaType(contentType)
		if mediaType != "text/plain" {
			continue
		}
		bytes, err := io.ReadAll(part.Body)
		if err != nil {
			return "", err
		}
		parts = append(parts, string(bytes))
	}
	return strings.Join(parts, "\n\n"), nil
}

func addressLine(header mail.Header, key string) string {
	addresses := addressList(header, key)
	if len(addresses) == 0 {
		return ""
	}
	return addresses[0]
}

func addressList(header mail.Header, key string) []string {
	addresses, err := header.AddressList(key)
	if err != nil {
		return nil
	}
	result := make([]string, 0, len(addresses))
	for _, address := range addresses {
		if address.Name != "" {
			result = append(result, address.Name+" <"+address.Address+">")
			continue
		}
		result = append(result, address.Address)
	}
	return result
}

func messageID(subject, from string, date time.Time, body string) string {
	sum := sha256.Sum256([]byte(subject + from + date.UTC().Format(time.RFC3339Nano) + body))
	return hex.EncodeToString(sum[:12])
}

func snippet(body string) string {
	compact := strings.Join(strings.Fields(body), " ")
	if len(compact) <= 180 {
		return compact
	}
	return compact[:177] + "..."
}

func fallback(value, replacement string) string {
	if strings.TrimSpace(value) == "" {
		return replacement
	}
	return strings.TrimSpace(value)
}
