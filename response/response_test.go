package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type xmlDocument struct {
	Value string `xml:"value"`
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()
	response := httptest.NewRecorder()
	WriteJSON(response, http.StatusCreated, map[string]string{"id": "example"})
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"id":"example"}`, response.Body.String())
}

func TestWriteXML(t *testing.T) {
	t.Parallel()
	response := httptest.NewRecorder()
	err := WriteXML(response, http.StatusOK, xmlDocument{Value: "example"})
	require.NoError(t, err)
	assert.Equal(t, "application/xml; charset=utf-8", response.Header().Get("Content-Type"))
	assert.Contains(t, response.Body.String(), "<value>example</value>")
}

func TestIsLocalPath(t *testing.T) {
	t.Parallel()
	assert.True(t, IsLocalPath("/pages/start?search=hello#section"))
	assert.False(t, IsLocalPath("https://evil.example"))
	assert.False(t, IsLocalPath("//evil.example"))
	assert.False(t, IsLocalPath("/\\evil.example"))
}

func TestStatusWriter(t *testing.T) {
	t.Parallel()
	response := httptest.NewRecorder()
	writer := NewStatusWriter(response)
	_, err := writer.Write([]byte("example"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, writer.Status)
	assert.Equal(t, "example", response.Body.String())
}
