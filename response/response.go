// Package response provides small, transport-neutral HTTP response helpers.
package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

// WriteJSON writes value as a JSON response.
func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}

// WriteXML writes value as indented XML with an XML declaration.
func WriteXML(w http.ResponseWriter, status int, value any) error {
	body, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(append(append([]byte{}, xml.Header...), append(body, '\n')...))
	return err
}

// IsLocalPath reports whether value is safe to use as a root-relative redirect.
func IsLocalPath(value string) bool {
	if !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") {
		return false
	}
	for _, character := range value {
		if character == '\\' || character <= 0x1f || character == 0x7f {
			return false
		}
	}
	return true
}
