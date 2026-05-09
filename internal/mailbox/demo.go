package mailbox

import "time"

func DemoMessages(now time.Time) []Message {
	return []Message{
		demoMessage(
			"demo-backend-1",
			"Re: Contract redlines before the investor call",
			"Maya North <maya@northstar.example>",
			[]string{"you@localhuman.dev"},
			now.Add(-18*time.Minute).UTC(),
			"Can you send the condensed version with only the liability deltas highlighted?",
			"Can you send the condensed version with only the liability deltas highlighted?\n\nThe long memo is useful, but the partner group is going to skim. I mostly need the parts that affect downside exposure, timeline, and approval authority.",
			"personal_reply",
			[]string{"imported", "personal_reply", "legal", "priority"},
		),
		demoMessage(
			"demo-backend-2",
			"Thursday launch checklist",
			"Ops <ops@atelier.example>",
			[]string{"you@localhuman.dev", "shipping@atelier.example"},
			now.Add(-2*time.Hour).UTC(),
			"The only unresolved items are the DNS cutover and the fallback status page copy.",
			"The only unresolved items are the DNS cutover and the fallback status page copy.\n\nEverything else is green. If you can confirm owner and timing by noon, we can keep Thursday without compressing QA.",
			"notification",
			[]string{"imported", "notification", "launch"},
		),
	}
}

func demoMessage(id, subject, from string, to []string, date time.Time, snippet, body, shape string, tags []string) Message {
	return Message{
		ID:          id,
		SourceID:    "demo",
		Subject:     subject,
		From:        from,
		To:          to,
		Date:        date,
		Snippet:     snippet,
		Body:        body,
		PrimaryBody: body,
		Shape:       shape,
		Tags:        tags,
		Confidence:  NewConfidence(0.9, "curated demo message"),
		FieldConfidence: map[string]Confidence{
			"subject": NewConfidence(0.95, "curated demo subject"),
			"from":    NewConfidence(0.95, "curated demo sender"),
			"date":    NewConfidence(0.95, "curated demo date"),
			"shape":   NewConfidence(0.85, "curated demo shape"),
			"body":    NewConfidence(0.9, "curated demo body"),
		},
		Warnings:    []ImportWarning{},
		Attachments: []Attachment{},
		Provenance: Provenance{
			SourceID:      "demo",
			SourceSHA256:  id,
			ParserVersion: ParserVersion,
			SchemaVersion: SchemaVersion,
			SizeBytes:     int64(len(body)),
		},
	}
}
