package tests

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

// TestAPIStatusCodes проверяет ошибки API по HTTP-кодам и JSON.
func TestAPIStatusCodes(t *testing.T) {
	cases := []struct {
		name, method, path, body string
		cookie                   bool
		status                   int
	}{
		{"missing task", "GET", "api/task?id=999999999", "", true, 404},
		{"delete missing task", "DELETE", "api/task?id=999999999", "", true, 404},
		{"invalid id", "GET", "api/task?id=bad", "", true, 400},
		{"invalid json", "POST", "api/task", "{", true, 400},
		{"invalid method", "PATCH", "api/task", "", true, 405},
		{"missing cookie", "GET", "api/tasks", "", false, 401},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			token := ""
			if item.cookie {
				token = Token
			}
			response, err := sendRequest(item.method, item.path, strings.NewReader(item.body), token)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			var result struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(body, &result); err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != item.status || result.Error == "" {
				t.Fatalf("unexpected response: %d %s", response.StatusCode, body)
			}
		})
	}
}
