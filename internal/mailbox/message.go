package mailbox

import "time"

type Message struct {
	ID      string    `json:"id"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	To      []string  `json:"to"`
	Date    time.Time `json:"date"`
	Snippet string    `json:"snippet"`
	Body    string    `json:"body"`
	Tags    []string  `json:"tags"`
}

type ImportResult struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
	Total    int `json:"total"`
}
