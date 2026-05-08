package search

import (
	"context"
	"testing"
	"time"

	"github.com/baditaflorin/localhuman-mail/internal/mailbox"
	"github.com/stretchr/testify/require"
)

func TestStoreImportsAndSearchesMessages(t *testing.T) {
	ctx := context.Background()
	store, err := Open(ctx, t.TempDir())
	require.NoError(t, err)
	defer store.Close()

	result, err := store.AddMessages(ctx, []mailbox.Message{
		{
			ID:      "one",
			Subject: "Launch checklist",
			From:    "ops@example.com",
			To:      []string{"you@example.com"},
			Date:    time.Now().UTC(),
			Snippet: "DNS cutover",
			Body:    "The DNS cutover is ready.",
			Tags:    []string{"launch"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1, result.Imported)

	messages, err := store.SearchMessages(ctx, "dns", 10)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	require.Equal(t, "one", messages[0].ID)
}
