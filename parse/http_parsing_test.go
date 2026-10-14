package parse

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	t.Parallel()

	t.Run("valid headers", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders([]string{"Content-Type=application/json", "Authorization=Bearer token"}, false)
		require.NoError(t, err)
		assert.Equal(t, http.Header{"Content-Type": {"application/json"}, "Authorization": {"Bearer token"}}, result)
	})

	t.Run("canonicalizes header names", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders([]string{"content-type=application/json"}, false)
		require.NoError(t, err)
		assert.Equal(t, http.Header{"Content-Type": {"application/json"}}, result)
	})

	t.Run("empty headers", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders(nil, false)
		require.NoError(t, err)
		assert.Equal(t, http.Header{}, result)
	})

	t.Run("malformed header", func(t *testing.T) {
		t.Parallel()
		_, err := ParseHeaders([]string{"AuthorizationBearer token"}, false)
		assert.EqualError(t, err, "invalid header format: AuthorizationBearer token")
	})

	t.Run("header with spaces", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders([]string{"  Content-Type = application/json  "}, false)
		require.NoError(t, err)
		assert.Equal(t, http.Header{"Content-Type": {"application/json"}}, result)
	})

	t.Run("header with empty key", func(t *testing.T) {
		t.Parallel()
		_, err := ParseHeaders([]string{"=value"}, false)
		assert.EqualError(t, err, "header key cannot be empty: =value")
	})

	t.Run("header with empty value", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders([]string{"X-Empty="}, false)
		require.NoError(t, err)
		assert.Equal(t, http.Header{"X-Empty": {""}}, result)
	})

	t.Run("preserves commas and equals in values", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders([]string{"Accept=text/html, application/json", "Authorization=Signature token=abc123"}, false)
		require.NoError(t, err)
		assert.Equal(t, "text/html, application/json", result.Get("Accept"))
		assert.Equal(t, "Signature token=abc123", result.Get("Authorization"))
	})

	t.Run("handles duplicate headers", func(t *testing.T) {
		t.Parallel()
		result, err := ParseHeaders([]string{"X-Test=one", "X-Test=two"}, true)
		require.NoError(t, err)
		assert.Equal(t, []string{"one", "two"}, result.Values("X-Test"))

		_, err = ParseHeaders([]string{"x-test=one", "X-Test=two"}, false)
		assert.EqualError(t, err, "duplicate header key found: X-Test")
	})
}

func TestParseStatusCodes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		input string
		want  []int
		err   string
	}{
		{name: "single code", input: "200", want: []int{200}},
		{name: "multiple codes", input: "200,404,500", want: []int{200, 404, 500}},
		{name: "ranges", input: "200-202,300-301,500", want: []int{200, 201, 202, 300, 301, 500}},
		{name: "invalid code", input: "abc", err: "invalid status code: abc"},
		{name: "invalid range", input: "200--202", err: "invalid status range: 200--202"},
		{name: "descending range", input: "201-200", err: "invalid status range: 201-200"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseStatusCodes(tc.input)
			if tc.err != "" {
				assert.EqualError(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
