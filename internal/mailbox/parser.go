package mailbox

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"regexp"
	"strings"
	"time"

	// Register charset decoders used by old exported mail such as ISO-8859-1.
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"golang.org/x/net/html"
)

type bodyCandidate struct {
	kind string
	text string
}

func ParseEML(reader io.Reader) (Message, error) {
	raw, err := readBounded(reader, MaxEMLBytes)
	if err != nil {
		return Message{}, err
	}
	sourceHash := sha256.Sum256(raw)
	sourceSHA := hex.EncodeToString(sourceHash[:])
	sourceID := "src_" + sourceSHA[:16]
	parseRaw, hadEnvelope := normalizeMboxEnvelope(raw)

	mailReader, err := mail.CreateReader(bytes.NewReader(parseRaw))
	if err != nil {
		return Message{}, ImportError{
			Kind:    "recoverable_input",
			What:    "Email could not be opened",
			Why:     "The file is not a readable RFC 5322/MIME message.",
			NowWhat: "Export the message as .eml again, then retry the import.",
		}
	}

	warnings := make([]ImportWarning, 0)
	if hadEnvelope {
		warnings = append(warnings, warning("info", "format", "Mbox envelope line detected and ignored.", "No action needed unless the sender looks wrong."))
	}
	fieldConfidence := make(map[string]Confidence)
	subject, subjectConfidence := parseSubject(mailReader.Header, &warnings)
	from, fromConfidence := parseAddressLine(mailReader.Header, "From", &warnings)
	to := addressList(mailReader.Header, "To")
	date, dateConfidence := parseDate(mailReader.Header, &warnings)
	fieldConfidence["subject"] = subjectConfidence
	fieldConfidence["from"] = fromConfidence
	fieldConfidence["date"] = dateConfidence

	candidates, attachments, calendar, partWarnings := collectParts(mailReader)
	warnings = append(warnings, partWarnings...)
	warnings = append(warnings, structuralWarnings(mailReader.Header, parseRaw)...)

	body, bodyKind := chooseBody(candidates)
	if strings.TrimSpace(body) == "" && calendar != nil {
		body = calendarText(calendar)
		bodyKind = "calendar"
	}
	if strings.TrimSpace(body) == "" {
		body = fallbackBodyFromRaw(parseRaw, boundaryOf(mailReader.Header.Get("Content-Type")))
		bodyKind = "partial"
		if body != "" {
			warnings = append(warnings, warning("warning", "body", "Readable text was recovered from a partial MIME message.", "Review the original message if content appears incomplete."))
		}
	}
	if strings.TrimSpace(body) == "" {
		return Message{}, ImportError{
			Kind:    "recoverable_input",
			What:    "Email has no readable body",
			Why:     "The importer could not find text, HTML, or calendar content to show.",
			NowWhat: "Try another export or keep the original message in your mail client for review.",
		}
	}

	primary, removedQuoted := primaryBody(body)
	if removedQuoted {
		warnings = append(warnings, warning("info", "body", "Quoted reply text was detected and de-emphasized.", "Review the readable body if the quoted thread matters."))
	}
	if bodyKind == "html" {
		warnings = append(warnings, warning("warning", "body", "Readable body was extracted from HTML because no better plaintext body was available.", "Review formatting-sensitive content before replying."))
	}
	if len(attachments) > 0 {
		warnings = append(warnings, warning("info", "attachments", fmt.Sprintf("%d attachment(s) detected.", len(attachments)), "Open the original email if attachment contents matter."))
	}
	if calendar != nil {
		warnings = append(warnings, warning("info", "calendar", "Calendar invite detected.", "Verify the event time before replying."))
	}

	shape, shapeConfidence := inferShape(mailReader.Header, subject, from, primary, bodyKind, attachments, calendar)
	fieldConfidence["shape"] = shapeConfidence
	bodyConfidence := bodyConfidence(bodyKind, warnings)
	fieldConfidence["body"] = bodyConfidence
	overall := overallConfidence(fieldConfidence, warnings)

	message := Message{
		ID:              messageID(mailReader.Header, sourceSHA),
		SourceID:        sourceID,
		Subject:         fallback(subject, "(no subject)"),
		From:            fallback(from, "unknown sender"),
		To:              to,
		Date:            date,
		Snippet:         snippet(primary),
		Body:            primary,
		PrimaryBody:     primary,
		Shape:           shape,
		Tags:            tagsForShape(shape),
		Confidence:      overall,
		FieldConfidence: fieldConfidence,
		Warnings:        warnings,
		Attachments:     attachments,
		Calendar:        calendar,
		Provenance: Provenance{
			SourceID:      sourceID,
			SourceSHA256:  sourceSHA,
			ParserVersion: ParserVersion,
			SchemaVersion: SchemaVersion,
			SizeBytes:     int64(len(raw)),
		},
	}
	return message, nil
}

func readBounded(reader io.Reader, maxBytes int64) ([]byte, error) {
	limited := io.LimitReader(reader, maxBytes+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read email: %w", err)
	}
	if int64(len(raw)) > maxBytes {
		return nil, ImportError{
			Kind:    "recoverable_input",
			What:    "Email is too large",
			Why:     fmt.Sprintf("The file is larger than the %dMB Phase 2 import budget.", maxBytes/(1024*1024)),
			NowWhat: "Import a smaller export or split the mailbox batch.",
		}
	}
	return raw, nil
}

func normalizeMboxEnvelope(raw []byte) ([]byte, bool) {
	if !bytes.HasPrefix(raw, []byte("From ")) {
		return raw, false
	}
	if index := bytes.IndexByte(raw, '\n'); index >= 0 && index+1 < len(raw) {
		return raw[index+1:], true
	}
	return raw, false
}

func parseSubject(header mail.Header, warnings *[]ImportWarning) (string, Confidence) {
	subject, err := header.Subject()
	if err != nil || strings.TrimSpace(subject) == "" {
		*warnings = append(*warnings, warning("warning", "subject", "Subject is missing or malformed.", "The message will import as '(no subject)'."))
		return "", NewConfidence(0.3, "subject header missing or malformed")
	}
	return normalizeHeader(subject), NewConfidence(0.95, "subject header parsed")
}

func fallbackBodyFromRaw(raw []byte, boundary string) string {
	text := strings.ReplaceAll(string(raw), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	body := text
	if index := strings.Index(text, "\n\n"); index >= 0 {
		body = text[index+2:]
	}
	if boundary != "" {
		marker := "--" + boundary
		if index := strings.Index(body, marker); index >= 0 {
			body = body[index+len(marker):]
		}
		body = strings.TrimLeft(body, "\n")
		if index := strings.Index(body, "\n\n"); index >= 0 {
			body = body[index+2:]
		}
		if index := strings.Index(body, marker); index >= 0 {
			body = body[:index]
		}
	}
	return normalizeText(body)
}

func parseDate(header mail.Header, warnings *[]ImportWarning) (time.Time, Confidence) {
	date, err := header.Date()
	if err != nil {
		*warnings = append(*warnings, warning("warning", "date", "Date header is missing or malformed.", "The message uses a deterministic unknown-date value instead of pretending it arrived now."))
		return time.Unix(0, 0).UTC(), NewConfidence(0.2, "date header missing or malformed")
	}
	return date.UTC(), NewConfidence(0.95, "date header parsed")
}

func parseAddressLine(header mail.Header, key string, warnings *[]ImportWarning) (string, Confidence) {
	addresses := addressList(header, key)
	if len(addresses) == 0 {
		*warnings = append(*warnings, warning("warning", strings.ToLower(key), key+" header is missing or malformed.", "Verify the sender in the original email."))
		return "", NewConfidence(0.2, key+" header missing or malformed")
	}
	return addresses[0], NewConfidence(0.9, key+" address parsed")
}

func addressList(header mail.Header, key string) []string {
	addresses, err := header.AddressList(key)
	if err != nil {
		return nil
	}
	result := make([]string, 0, len(addresses))
	for _, address := range addresses {
		if address.Name != "" {
			result = append(result, normalizeHeader(address.Name)+" <"+address.Address+">")
			continue
		}
		result = append(result, address.Address)
	}
	return result
}

func collectParts(mailReader *mail.Reader) ([]bodyCandidate, []Attachment, *CalendarEvent, []ImportWarning) {
	candidates := make([]bodyCandidate, 0)
	attachments := make([]Attachment, 0)
	warnings := make([]ImportWarning, 0)
	var calendar *CalendarEvent

	for {
		part, err := mailReader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			warnings = append(warnings, warning("warning", "mime", "Email appears partial or has a malformed MIME boundary.", "Retry the export if content seems missing."))
			break
		}

		switch header := part.Header.(type) {
		case *mail.InlineHeader:
			contentType := contentTypeOf(header)
			text, readErr := readPartText(part.Body)
			if readErr != nil {
				warnings = append(warnings, warning("warning", "body", "A message part could not be read completely.", "Retry the export if content seems missing."))
				continue
			}
			switch {
			case strings.HasPrefix(contentType, "text/plain"):
				candidates = append(candidates, bodyCandidate{kind: "plain", text: text})
			case strings.HasPrefix(contentType, "text/html"):
				candidates = append(candidates, bodyCandidate{kind: "html", text: htmlToText(text)})
			case strings.HasPrefix(contentType, "text/calendar"):
				parsed := parseCalendar(text)
				calendar = &parsed
				candidates = append(candidates, bodyCandidate{kind: "calendar", text: calendarText(&parsed)})
			}
		case *mail.AttachmentHeader:
			contentType := contentTypeOf(header)
			fileName, _ := header.Filename()
			size, _ := countPart(part.Body)
			attachments = append(attachments, Attachment{
				FileName:    fallback(fileName, "(unnamed attachment)"),
				ContentType: fallback(contentType, "application/octet-stream"),
				SizeBytes:   size,
			})
		}
	}
	return candidates, attachments, calendar, warnings
}

func contentTypeOf(header interface {
	ContentType() (string, map[string]string, error)
}) string {
	contentType, _, err := header.ContentType()
	if err != nil {
		return "text/plain"
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return strings.ToLower(contentType)
	}
	return strings.ToLower(mediaType)
}

func readPartText(reader io.Reader) (string, error) {
	raw, err := io.ReadAll(io.LimitReader(reader, MaxIndexedBytes+1))
	if err != nil {
		return "", err
	}
	if len(raw) > MaxIndexedBytes {
		raw = raw[:MaxIndexedBytes]
	}
	return string(raw), nil
}

func countPart(reader io.Reader) (int64, error) {
	return io.Copy(io.Discard, io.LimitReader(reader, MaxIndexedBytes+1))
}

func chooseBody(candidates []bodyCandidate) (string, string) {
	for _, candidate := range candidates {
		if candidate.kind == "plain" && strings.TrimSpace(candidate.text) != "" {
			return candidate.text, candidate.kind
		}
	}
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.text) != "" {
			return candidate.text, candidate.kind
		}
	}
	return "", ""
}

func structuralWarnings(header mail.Header, raw []byte) []ImportWarning {
	warnings := make([]ImportWarning, 0)
	contentType := header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "charset=") &&
		!strings.Contains(strings.ToLower(contentType), "utf-8") &&
		!strings.Contains(strings.ToLower(contentType), "us-ascii") {
		warnings = append(warnings, warning("info", "charset", "Non-UTF charset detected and normalized.", "Verify unusual characters before replying."))
	}
	if boundary := boundaryOf(contentType); boundary != "" && !bytes.Contains(raw, []byte("--"+boundary+"--")) {
		warnings = append(warnings, warning("warning", "mime", "MIME message appears partial; closing boundary is missing.", "Retry the export if the message looks incomplete."))
	}
	if header.Get("List-Id") != "" {
		warnings = append(warnings, warning("info", "list", "Mailing list headers detected.", "Use list context when replying."))
	}
	if header.Get("List-Unsubscribe") != "" {
		warnings = append(warnings, warning("info", "unsubscribe", "Unsubscribe header detected.", "This is probably a list or newsletter."))
	}
	return warnings
}

func boundaryOf(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}
	return params["boundary"]
}

func htmlToText(value string) string {
	node, err := html.Parse(strings.NewReader(value))
	if err != nil {
		return value
	}
	var parts []string
	var walk func(*html.Node)
	walk = func(current *html.Node) {
		if current.Type == html.ElementNode && (current.Data == "script" || current.Data == "style") {
			return
		}
		if current.Type == html.TextNode {
			parts = append(parts, current.Data)
		}
		for child := current.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return strings.Join(parts, " ")
}

func normalizeText(value string) string {
	value = strings.TrimPrefix(value, "\ufeff")
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.ReplaceAll(value, "\u00a0", " ")
	value = strings.ReplaceAll(value, "\t", " ")
	return strings.Join(strings.Fields(value), " ")
}

func normalizeHeader(value string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(value, "\u00a0", " ")), " ")
}

func primaryBody(value string) (string, bool) {
	lines := strings.Split(value, "\n")
	if len(lines) == 1 {
		return value, false
	}
	filtered := make([]string, 0, len(lines))
	removed := false
	replyIntro := regexp.MustCompile(`(?i)^on .+ wrote:$`)
	namedQuote := regexp.MustCompile(`^[A-Za-z0-9_. -]{1,30}>`)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ">") || replyIntro.MatchString(trimmed) || namedQuote.MatchString(trimmed) {
			removed = true
			continue
		}
		filtered = append(filtered, line)
	}
	return normalizeText(strings.Join(filtered, "\n")), removed
}

func parseCalendar(value string) CalendarEvent {
	lines := strings.Split(value, "\n")
	event := CalendarEvent{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "SUMMARY:"):
			event.Summary = strings.TrimPrefix(line, "SUMMARY:")
		case strings.HasPrefix(line, "LOCATION:"):
			event.Location = strings.TrimPrefix(line, "LOCATION:")
		case strings.HasPrefix(line, "DTSTART:"):
			event.Start = strings.TrimPrefix(line, "DTSTART:")
		case strings.HasPrefix(line, "DTEND:"):
			event.End = strings.TrimPrefix(line, "DTEND:")
		}
	}
	return event
}

func calendarText(event *CalendarEvent) string {
	if event == nil {
		return ""
	}
	return normalizeText(strings.Join([]string{event.Summary, event.Location, event.Start, event.End}, " "))
}

func inferShape(header mail.Header, subject, from, body, bodyKind string, attachments []Attachment, calendar *CalendarEvent) (string, Confidence) {
	joined := strings.ToLower(subject + " " + from + " " + body)
	switch {
	case calendar != nil:
		return "calendar_invite", NewConfidence(0.92, "text/calendar part detected")
	case len(attachments) > 0:
		return "attachment_only", NewConfidence(0.72, "attachment metadata detected")
	case strings.Contains(strings.ToLower(header.Get("Content-Type")), "multipart/mixed"):
		return "attachment_only", NewConfidence(0.45, "multipart/mixed envelope suggests attachments but content is partial")
	case header.Get("List-Id") != "":
		return "mailing_list", NewConfidence(0.9, "List-Id header detected")
	case header.Get("List-Unsubscribe") != "" || strings.Contains(joined, "newsletter") || strings.Contains(strings.ToLower(from), "newsletter"):
		return "newsletter", NewConfidence(0.82, "newsletter/unsubscribe evidence detected")
	case bodyKind == "html" && strings.Contains(joined, "unsubscribe"):
		return "newsletter", NewConfidence(0.72, "HTML unsubscribe evidence detected")
	case strings.Contains(strings.ToLower(from), "rssfeeds") || strings.Contains(joined, "notification"):
		return "notification", NewConfidence(0.75, "notification sender/body evidence detected")
	case strings.Contains(joined, "invoice") || strings.Contains(joined, "receipt"):
		return "receipt_invoice", NewConfidence(0.75, "receipt/invoice language detected")
	case header.Get("In-Reply-To") != "" || header.Get("References") != "":
		return "personal_reply", NewConfidence(0.82, "thread headers detected")
	case bodyKind == "html":
		return "newsletter", NewConfidence(0.65, "HTML-only message resembles newsletter")
	default:
		return UnknownShape, NewConfidence(0.55, "no strong domain shape detected")
	}
}

func bodyConfidence(kind string, warnings []ImportWarning) Confidence {
	score := 0.9
	reasons := []string{"readable body extracted"}
	if kind == "html" {
		score = 0.68
		reasons = append(reasons, "body came from HTML fallback")
	}
	for _, warning := range warnings {
		if warning.Field == "mime" {
			score -= 0.25
			reasons = append(reasons, "MIME was partial or malformed")
			break
		}
	}
	if score < 0.1 {
		score = 0.1
	}
	return NewConfidence(score, reasons...)
}

func overallConfidence(fields map[string]Confidence, warnings []ImportWarning) Confidence {
	total := 0.0
	count := 0
	for _, confidence := range fields {
		total += confidence.Score
		count++
	}
	if count == 0 {
		return NewConfidence(0.3, "no fields were inferred")
	}
	score := total / float64(count)
	for _, warning := range warnings {
		if warning.Severity == "warning" {
			score -= 0.05
		}
	}
	if score < 0.1 {
		score = 0.1
	}
	return NewConfidence(score, "combined field confidence and parser warnings")
}

func messageID(header mail.Header, sourceSHA string) string {
	messageID := strings.Trim(header.Get("Message-ID"), "<> ")
	if messageID == "" {
		return "msg_" + sourceSHA[:24]
	}
	sum := sha256.Sum256([]byte(messageID + ":" + sourceSHA))
	return "msg_" + hex.EncodeToString(sum[:12])
}

func tagsForShape(shape string) []string {
	if shape == UnknownShape {
		return []string{"imported"}
	}
	return []string{"imported", shape}
}

func snippet(body string) string {
	compact := normalizeText(body)
	if len(compact) <= 180 {
		return compact
	}
	return compact[:177] + "..."
}

func warning(severity, field, message, nextStep string) ImportWarning {
	return ImportWarning{Severity: severity, Field: field, Message: message, NextStep: nextStep}
}

func fallback(value, replacement string) string {
	if strings.TrimSpace(value) == "" {
		return replacement
	}
	return strings.TrimSpace(value)
}
