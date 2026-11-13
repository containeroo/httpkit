package parse

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultStatusCodeMin = 100
	defaultStatusCodeMax = 599
)

// StatusCodeOption configures the accepted HTTP status code range.
type StatusCodeOption func(*statusCodeBounds)

// statusCodeBounds limits parsed HTTP status codes and range expansion.
type statusCodeBounds struct {
	min int
	max int
}

// WithStatusCodeRange sets inclusive minimum and maximum HTTP status codes.
// Both values must be within 100–599, and min must not exceed max.
// Invalid bounds cause ParseStatusCodes to return an error.
func WithStatusCodeRange(min, max int) StatusCodeOption {
	return func(bounds *statusCodeBounds) {
		bounds.min = min
		bounds.max = max
	}
}

// ParseStatusCodes parses comma-separated HTTP status codes and inclusive
// ranges, such as "200-204,429,500-503". By default, codes must be in the
// HTTP range 100–599. WithStatusCodeRange can restrict the accepted range.
// Endpoints are checked before expansion to avoid unbounded allocations.
// Input order and duplicates are preserved. Invalid input returns an error.
func ParseStatusCodes(value string, options ...StatusCodeOption) ([]int, error) {
	bounds := statusCodeBounds{
		min: defaultStatusCodeMin,
		max: defaultStatusCodeMax,
	}

	for _, option := range options {
		if option != nil {
			option(&bounds)
		}
	}

	if err := bounds.validate(); err != nil {
		return nil, err
	}

	return parseStatusCodeRanges(value, bounds)
}

// parseStatusCodeRanges expands validated single codes and ranges in input order.
func parseStatusCodeRanges(value string, bounds statusCodeBounds) ([]int, error) {
	var codes []int
	for _, entry := range strings.Split(value, ",") {
		start, end, err := parseStatusCodeEntry(strings.TrimSpace(entry), bounds)
		if err != nil {
			return nil, err
		}

		for code := start; code <= end; code++ {
			codes = append(codes, code)
		}
	}

	return codes, nil
}

// validate checks option boundaries before any status codes are parsed.
func (b statusCodeBounds) validate() error {
	if b.min > b.max {
		return fmt.Errorf("invalid HTTP status code range: minimum %d exceeds maximum %d", b.min, b.max)
	}
	if b.min < defaultStatusCodeMin || b.max > defaultStatusCodeMax {
		return fmt.Errorf("invalid HTTP status code range %d-%d: bounds must be between %d and %d", b.min, b.max, defaultStatusCodeMin, defaultStatusCodeMax)
	}
	return nil
}

// parseStatusCodeEntry parses and validates endpoints before range expansion.
func parseStatusCodeEntry(entry string, bounds statusCodeBounds) (int, int, error) {
	if entry == "" {
		return 0, 0, fmt.Errorf("invalid status code: empty entry")
	}

	startText, endText, isRange := strings.Cut(entry, "-")
	if !isRange {
		code, err := strconv.Atoi(entry)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid status code: %s", entry)
		}
		if err := bounds.validateCode(code); err != nil {
			return 0, 0, err
		}
		return code, code, nil
	}

	if strings.Contains(endText, "-") {
		return 0, 0, fmt.Errorf("invalid status range: %s", entry)
	}

	start, startErr := strconv.Atoi(strings.TrimSpace(startText))
	end, endErr := strconv.Atoi(strings.TrimSpace(endText))
	if startErr != nil || endErr != nil || start > end {
		return 0, 0, fmt.Errorf("invalid status range: %s", entry)
	}

	if err := bounds.validateCode(start); err != nil {
		return 0, 0, err
	}
	if err := bounds.validateCode(end); err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

// validateCode checks a parsed status against the configured boundaries.
func (b statusCodeBounds) validateCode(code int) error {
	if code < b.min || code > b.max {
		return fmt.Errorf("HTTP status code %d is outside the allowed range %d-%d", code, b.min, b.max)
	}
	return nil
}
