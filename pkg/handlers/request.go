package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

const maxRequestBytes = 1 << 20

// readJSONRequest читает JSON-объект размером до одного мегабайта.
func readJSONRequest(w http.ResponseWriter, r *http.Request, value any) bool {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" {
		writeAPIError(w, http.StatusUnsupportedMediaType, "content type must be application/json")
		return false
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxRequestBytes))
	if err != nil {
		var limitError *http.MaxBytesError
		if errors.As(err, &limitError) {
			writeAPIError(w, http.StatusRequestEntityTooLarge, "request body is too large")
		} else {
			writeAPIError(w, http.StatusBadRequest, "could not read request body")
		}
		return false
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 || body[0] != '{' || json.Unmarshal(body, value) != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid json")
		return false
	}
	return true
}
