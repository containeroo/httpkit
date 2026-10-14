// Package parse parses common HTTP configuration values.
package parse

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ParseStatusCodes parses a comma-separated string of HTTP status codes and ranges
// into a slice of integers. It accepts individual codes (for example, "200") and
// inclusive ranges (for example, "200-204").
func ParseStatusCodes(statusRanges string) ([]int, error) {
	var statusCodes []int

	ranges := strings.Split(statusRanges, ",")
	for _, r := range ranges {
		trimmed := strings.TrimSpace(r)

		if !strings.Contains(trimmed, "-") {
			code, err := strconv.Atoi(trimmed)
			if err != nil {
				return nil, fmt.Errorf("invalid status code: %s", trimmed)
			}
			statusCodes = append(statusCodes, code)
			continue
		}

		parts := strings.Split(trimmed, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid status range: %s", trimmed)
		}

		start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil || start > end {
			return nil, fmt.Errorf("invalid status range: %s", trimmed)
		}

		for i := start; i <= end; i++ {
			statusCodes = append(statusCodes, i)
		}
	}

	return statusCodes, nil
}

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
