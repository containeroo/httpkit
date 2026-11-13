// Package parse parses common HTTP configuration values.
package parse

import (
	"fmt"
	"net/http"
	"strings"
)

// ParseHeaders parses Key=Value entries into an http.Header. Values may be empty,
// but keys must not be. Header names are canonicalized. Duplicate values are
// preserved only when allowDuplicates is true.
func ParseHeaders(headers []string, allowDuplicates bool) (http.Header, error) {
	result := make(http.Header, len(headers))

	for _, raw := range headers {
		key, value, ok := strings.Cut(raw, "=")
		if !ok {
			return nil, fmt.Errorf("invalid header format: %s", raw)
		}

		key = http.CanonicalHeaderKey(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		if key == "" {
			return nil, fmt.Errorf("header key cannot be empty: %s", raw)
		}

		if _, exists := result[key]; exists && !allowDuplicates {
			return nil, fmt.Errorf("duplicate header key found: %s", key)
		}

		if allowDuplicates {
			result.Add(key, value)
			continue
		}
		result.Set(key, value)
	}

	return result, nil
}
