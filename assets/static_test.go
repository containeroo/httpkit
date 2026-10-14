package assets

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestStaticCachingAndCompression(t *testing.T) {
	data := bytes.Repeat([]byte("hello, world!"), 200)
	handler := Static(StaticConfig{Files: fstest.MapFS{
		"assets/app.js": {Data: data}, "sw.js": {Data: []byte("service worker")},
	}, ImmutablePrefix: "/assets/", NoCacheFiles: []string{"sw.js"}})
	req := httptest.NewRequest("GET", "/assets/app.js", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	result := httptest.NewRecorder()
	handler.ServeHTTP(result, req)
	if result.Code != 200 || result.Header().Get("Content-Encoding") != "gzip" || result.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" {
		t.Fatalf("bad response: %#v", result)
	}
	reader, err := gzip.NewReader(bytes.NewReader(result.Body.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, data) {
		t.Fatal("corrupt gzip response")
	}
	result = httptest.NewRecorder()
	handler.ServeHTTP(result, httptest.NewRequest(http.MethodGet, "/sw.js", nil))
	if result.Code != 200 || result.Header().Get("Cache-Control") != "no-cache" {
		t.Fatal("service worker caching")
	}
}

func TestAcceptsGzip(t *testing.T) {
	for _, tc := range []struct {
		header string
		want   bool
	}{
		{"gzip", true}, {"gzip;q=0", false}, {"br, *;q=0.2", true}, {"gzip;q=0, *;q=1", false}, {"br", false},
	} {
		if got := AcceptsGzip(tc.header); got != tc.want {
			t.Errorf("%q: %t want %t", tc.header, got, tc.want)
		}
	}
}
