package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStatusCodes(t *testing.T) {
	t.Parallel()

	t.Run("single status code", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200")

		require.NoError(t, err)
		assert.Equal(t, []int{200}, codes)
	})

	t.Run("multiple status codes", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200,404,500")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 404, 500}, codes)
	})

	t.Run("inclusive ranges", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-204,429,500-503")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 201, 202, 203, 204, 429, 500, 501, 502, 503}, codes)
	})

	t.Run("whitespace around entries and endpoints", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes(" 200 - 202 ,  429 , 500 - 501 ")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 201, 202, 429, 500, 501}, codes)
	})

	t.Run("default lower and upper boundaries", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("100,599")

		require.NoError(t, err)
		assert.Equal(t, []int{100, 599}, codes)
	})

	t.Run("all default HTTP status codes", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("100-599")

		require.NoError(t, err)
		assert.Len(t, codes, 500)
		assert.Equal(t, 100, codes[0])
		assert.Equal(t, 599, codes[len(codes)-1])
	})

	t.Run("preserves duplicate values and input order", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("429,200-201,200,429")

		require.NoError(t, err)
		assert.Equal(t, []int{429, 200, 201, 200, 429}, codes)
	})

	t.Run("single value ranges", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-200,500-500")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 500}, codes)
	})
}

func TestParseStatusCodesInvalidInput(t *testing.T) {
	t.Parallel()

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("")

		require.EqualError(t, err, "invalid status code: empty entry")
		assert.Nil(t, codes)
	})

	t.Run("whitespace-only input", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("  ")

		require.EqualError(t, err, "invalid status code: empty entry")
		assert.Nil(t, codes)
	})

	t.Run("empty entry in middle", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200,,201")

		require.ErrorContains(t, err, "empty entry")
		assert.Nil(t, codes)
	})

	t.Run("trailing comma", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200,")

		require.ErrorContains(t, err, "empty entry")
		assert.Nil(t, codes)
	})

	t.Run("leading comma", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes(",200")

		require.ErrorContains(t, err, "empty entry")
		assert.Nil(t, codes)
	})

	t.Run("non-numeric status", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200,abc")

		require.ErrorContains(t, err, `invalid status code: abc`)
		assert.Nil(t, codes)
	})

	t.Run("default status too low", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("99")

		require.EqualError(t, err, "HTTP status code 99 is outside the allowed range 100-599")
		assert.Nil(t, codes)
	})

	t.Run("default status too high", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("600")

		require.EqualError(t, err, "HTTP status code 600 is outside the allowed range 100-599")
		assert.Nil(t, codes)
	})

	t.Run("negative status", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("-1")

		require.ErrorContains(t, err, "invalid status range")
		assert.Nil(t, codes)
	})

	t.Run("descending range", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("204-200")

		require.EqualError(t, err, `invalid status range: 204-200`)
		assert.Nil(t, codes)
	})

	t.Run("range missing end", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-")

		require.ErrorContains(t, err, "invalid status range")
		assert.Nil(t, codes)
	})

	t.Run("range missing start", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("-200")

		require.ErrorContains(t, err, "invalid status range")
		assert.Nil(t, codes)
	})

	t.Run("range with extra separator", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200--204")

		require.EqualError(t, err, `invalid status range: 200--204`)
		assert.Nil(t, codes)
	})

	t.Run("range with non-numeric start", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("abc-204")

		require.ErrorContains(t, err, "invalid status range")
		assert.Nil(t, codes)
	})

	t.Run("range with non-numeric end", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-def")

		require.ErrorContains(t, err, "invalid status range")
		assert.Nil(t, codes)
	})

	t.Run("range starts below default minimum", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("99-200")

		require.EqualError(t, err, "HTTP status code 99 is outside the allowed range 100-599")
		assert.Nil(t, codes)
	})

	t.Run("range ends above default maximum", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-600")

		require.EqualError(t, err, "HTTP status code 600 is outside the allowed range 100-599")
		assert.Nil(t, codes)
	})

	t.Run("huge range rejected before expansion", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-999999999")

		require.EqualError(t, err, "HTTP status code 999999999 is outside the allowed range 100-599")
		assert.Nil(t, codes)
	})

	t.Run("range endpoint integer overflow", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-999999999999999999999999999999999999999")

		require.ErrorContains(t, err, "invalid status range")
		assert.Nil(t, codes)
	})

	t.Run("single status integer overflow", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("999999999999999999999999999999")

		require.ErrorContains(t, err, "invalid status code")
		assert.Nil(t, codes)
	})
}

func TestParseStatusCodesOptions(t *testing.T) {
	t.Parallel()

	t.Run("custom range accepts codes and ranges", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-204,429,500-503", WithStatusCodeRange(200, 599))

		require.NoError(t, err)
		assert.Equal(t, []int{200, 201, 202, 203, 204, 429, 500, 501, 502, 503}, codes)
	})

	t.Run("custom range includes both endpoints", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200,299", WithStatusCodeRange(200, 299))

		require.NoError(t, err)
		assert.Equal(t, []int{200, 299}, codes)
	})

	t.Run("custom range rejects code below minimum", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("199,200", WithStatusCodeRange(200, 599))

		require.EqualError(t, err, "HTTP status code 199 is outside the allowed range 200-599")
		assert.Nil(t, codes)
	})

	t.Run("custom range rejects code above maximum", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200,400", WithStatusCodeRange(200, 399))

		require.EqualError(t, err, "HTTP status code 400 is outside the allowed range 200-399")
		assert.Nil(t, codes)
	})

	t.Run("custom range rejects lower range endpoint", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("199-202", WithStatusCodeRange(200, 599))

		require.EqualError(t, err, "HTTP status code 199 is outside the allowed range 200-599")
		assert.Nil(t, codes)
	})

	t.Run("custom range rejects upper range endpoint", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-400", WithStatusCodeRange(200, 399))

		require.EqualError(t, err, "HTTP status code 400 is outside the allowed range 200-399")
		assert.Nil(t, codes)
	})

	t.Run("equal option boundaries accept only that code", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("204-204", WithStatusCodeRange(204, 204))

		require.NoError(t, err)
		assert.Equal(t, []int{204}, codes)
	})

	t.Run("minimum exceeds maximum", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-204", WithStatusCodeRange(599, 200))

		require.EqualError(t, err, "invalid HTTP status code range: minimum 599 exceeds maximum 200")
		assert.Nil(t, codes)
	})

	t.Run("option minimum below valid HTTP range", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200", WithStatusCodeRange(99, 599))

		require.EqualError(t, err, "invalid HTTP status code range 99-599: bounds must be between 100 and 599")
		assert.Nil(t, codes)
	})

	t.Run("option maximum above valid HTTP range", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200", WithStatusCodeRange(100, 600))

		require.EqualError(t, err, "invalid HTTP status code range 100-600: bounds must be between 100 and 599")
		assert.Nil(t, codes)
	})

	t.Run("option zero bounds are rejected", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200", WithStatusCodeRange(0, 0))

		require.EqualError(t, err, "invalid HTTP status code range 0-0: bounds must be between 100 and 599")
		assert.Nil(t, codes)
	})

	t.Run("invalid option is checked before input", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("not-a-status", WithStatusCodeRange(599, 200))

		require.EqualError(t, err, "invalid HTTP status code range: minimum 599 exceeds maximum 200")
		assert.Nil(t, codes)
	})

	t.Run("nil option is ignored", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("204", nil)

		require.NoError(t, err)
		assert.Equal(t, []int{204}, codes)
	})

	t.Run("later options override earlier options", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("500", WithStatusCodeRange(200, 299), WithStatusCodeRange(200, 599))

		require.NoError(t, err)
		assert.Equal(t, []int{500}, codes)
	})
}

func TestParseStatusCodesChangedBehavior(t *testing.T) {
	t.Parallel()

	t.Run("out-of-range codes are rejected by the consolidated parser", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("99,600")

		require.EqualError(t, err, "HTTP status code 99 is outside the allowed range 100-599")
		assert.Nil(t, codes)
	})

	t.Run("original valid range parsing still works", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200-202,500")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 201, 202, 500}, codes)
	})
}

func TestParseStatusCodesSyntaxCompatibility(t *testing.T) {
	t.Parallel()

	t.Run("invalid status code retains error format", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("abc")

		require.EqualError(t, err, "invalid status code: abc")
		assert.Nil(t, codes)
	})

	t.Run("invalid range retains error format", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("200--202")

		require.EqualError(t, err, "invalid status range: 200--202")
		assert.Nil(t, codes)
	})

	t.Run("descending range retains error format", func(t *testing.T) {
		t.Parallel()

		codes, err := ParseStatusCodes("201-200")

		require.EqualError(t, err, "invalid status range: 201-200")
		assert.Nil(t, codes)
	})
}
