package handlers

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vlonebara/todo-list/pkg/model"
)

// TestWriteTaskError проверяет коды ответов и скрытие внутренних ошибок.
func TestWriteTaskError(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{"invalid input", fmt.Errorf("%w: title is required", model.ErrInvalidTask), 400, "title is required"},
		{"missing task", model.ErrTaskNotFound, 404, "task not found"},
		{"database failure", errors.New("database path and private details"), 500, "internal server error"},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			writeTaskError(response, item.err)
			if response.Code != item.status || !strings.Contains(response.Body.String(), item.message) {
				t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "private details") {
				t.Fatal("internal error was exposed")
			}
		})
	}
}
