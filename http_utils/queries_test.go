package httputils

import "testing"

func TestPageQueries(t *testing.T) {
	t.Parallel()

	limit, cursor := 50, "abc"

	if got := PageQueries(nil, nil).Encode(); got != "" {
		t.Errorf("PageQueries(nil, nil) = %q, want empty", got)
	}

	if got := PageQueries(&limit, &cursor).Encode(); got != "cursor=abc&limit=50" {
		t.Errorf("PageQueries() = %q", got)
	}
}
