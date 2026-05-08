package mailbox

import "time"

func DemoMessages(now time.Time) []Message {
	return []Message{
		{
			ID:      "demo-backend-1",
			Subject: "Re: Contract redlines before the investor call",
			From:    "Maya North <maya@northstar.example>",
			To:      []string{"you@localhuman.dev"},
			Date:    now.Add(-18 * time.Minute).UTC(),
			Snippet: "Can you send the condensed version with only the liability deltas highlighted?",
			Body:    "Can you send the condensed version with only the liability deltas highlighted?\n\nThe long memo is useful, but the partner group is going to skim. I mostly need the parts that affect downside exposure, timeline, and approval authority.",
			Tags:    []string{"legal", "priority"},
		},
		{
			ID:      "demo-backend-2",
			Subject: "Thursday launch checklist",
			From:    "Ops <ops@atelier.example>",
			To:      []string{"you@localhuman.dev", "shipping@atelier.example"},
			Date:    now.Add(-2 * time.Hour).UTC(),
			Snippet: "The only unresolved items are the DNS cutover and the fallback status page copy.",
			Body:    "The only unresolved items are the DNS cutover and the fallback status page copy.\n\nEverything else is green. If you can confirm owner and timing by noon, we can keep Thursday without compressing QA.",
			Tags:    []string{"launch"},
		},
	}
}
