// Package assets serves embedded files with configurable caching and gzip compression.
package assets

import (
	"bytes"
	"compress/gzip"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StaticConfig configures a read-only embedded file server.
type StaticConfig struct {
	// Files is a filesystem rooted at the public asset directory.
	Files fs.FS
	// ImmutablePrefix marks fingerprinted build assets, for example /assets/.
	ImmutablePrefix string
	// NoCacheFiles are asset paths whose content must be revalidated.
	NoCacheFiles []string
	// MinCompressSize is the minimum size for gzip, defaulting to 1024 bytes.
	MinCompressSize int
}

// Static serves embedded files with immutable asset caching and negotiated gzip.
func Static(config StaticConfig) http.Handler {
	minimum := config.MinCompressSize
	if minimum <= 0 {
		minimum = 1024
	}
	h := &staticHandler{config: config, files: http.FileServer(http.FS(config.Files)), minimum: minimum}
	return h
}

type staticHandler struct {
	config     StaticConfig
	files      http.Handler
	minimum    int
	compressed sync.Map
}

// ServeHTTP serves an asset, preserving file modtime and content type.
func (h *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	if !fs.ValidPath(name) || name == "." {
		http.NotFound(w, r)
		return
	}
	cache := "public, max-age=86400"
	if h.config.ImmutablePrefix != "" && strings.HasPrefix(r.URL.Path, h.config.ImmutablePrefix) {
		cache = "public, max-age=31536000, immutable"
	}
	for _, special := range h.config.NoCacheFiles {
		if name == strings.TrimPrefix(special, "/") {
			cache = "no-cache"
			break
		}
	}
	w.Header().Set("Cache-Control", cache)
	if !compressible(name) {
		h.files.ServeHTTP(w, r)
		return
	}
	w.Header().Add("Vary", "Accept-Encoding")
	if !AcceptsGzip(r.Header.Get("Accept-Encoding")) {
		h.files.ServeHTTP(w, r)
		return
	}
	data, err := fs.ReadFile(h.config.Files, name)
	if err != nil || len(data) < h.minimum {
		h.files.ServeHTTP(w, r)
		return
	}
	encoded, err := h.compress(name, data)
	if err != nil {
		h.files.ServeHTTP(w, r)
		return
	}
	if contentType := mime.TypeByExtension(path.Ext(name)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Content-Encoding", "gzip")
	modified := time.Time{}
	if stat, err := fs.Stat(h.config.Files, name); err == nil {
		modified = stat.ModTime()
	}
	http.ServeContent(w, r, path.Base(name), modified, bytes.NewReader(encoded))
}

// compress returns a cached gzip representation for one asset.
func (h *staticHandler) compress(name string, data []byte) ([]byte, error) {
	if stored, ok := h.compressed.Load(name); ok {
		return stored.([]byte), nil
	}
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	actual, _ := h.compressed.LoadOrStore(name, buffer.Bytes())
	return actual.([]byte), nil
}

// compressible reports whether an asset is text-like.
func compressible(name string) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".css", ".html", ".js", ".json", ".svg", ".txt":
		return true
	default:
		return false
	}
}

// AcceptsGzip reports whether Accept-Encoding permits gzip (or an eligible wildcard).
func AcceptsGzip(value string) bool {
	wildcard := false
	for _, entry := range strings.Split(value, ",") {
		parts := strings.Split(strings.TrimSpace(entry), ";")
		encoding := strings.ToLower(strings.TrimSpace(parts[0]))
		quality := 1.0
		for _, parameter := range parts[1:] {
			kv := strings.SplitN(strings.TrimSpace(parameter), "=", 2)
			if len(kv) != 2 || !strings.EqualFold(kv[0], "q") {
				continue
			}
			if parsed, err := strconv.ParseFloat(kv[1], 64); err == nil {
				quality = parsed
			}
		}
		if encoding == "gzip" {
			return quality > 0
		}
		if encoding == "*" && quality > 0 {
			wildcard = true
		}
	}
	return wildcard
}
