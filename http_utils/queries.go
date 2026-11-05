package httputils

import (
	"net/url"
	"strconv"
)

// PageQueries builds the cursor and limit query parameters of a list call; nil values are left
// out so the API applies its defaults.
func PageQueries(limit *int, cursor *string) url.Values {
	queries := url.Values{}
	if cursor != nil {
		queries.Set("cursor", *cursor)
	}

	if limit != nil {
		queries.Set("limit", strconv.Itoa(*limit))
	}

	return queries
}
