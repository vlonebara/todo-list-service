package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

// TestHistoryAndExport проверяет новые маршруты через работающий HTTP-сервер.
func TestHistoryAndExport(t *testing.T) {
	result, err := postJSON("api/task", map[string]any{"date": "20300101", "title": "history export test"}, http.MethodPost)
	if err != nil {
		t.Fatal(err)
	}
	id, ok := result["id"].(string)
	if !ok || id == "" {
		t.Fatal("task was not created")
	}
	body, err := requestJSON("api/tasks/export", nil, http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}
	var exported struct {
		Tasks []map[string]string `json:"tasks"`
	}
	if err := json.Unmarshal(body, &exported); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, task := range exported.Tasks {
		if task["id"] == id {
			found = true
		}
	}
	if !found {
		t.Fatal("created task was not exported")
	}
	body, err = requestJSON("api/history", nil, http.MethodGet)
	if err != nil {
		t.Fatal(err)
	}
	var history struct {
		History []struct {
			Action string            `json:"action"`
			After  map[string]string `json:"after"`
		} `json:"history"`
	}
	if err := json.Unmarshal(body, &history); err != nil {
		t.Fatal(err)
	}
	if len(history.History) == 0 || history.History[0].After["id"] != id {
		t.Fatal("created task was not recorded")
	}
	if _, err := postJSON("api/task/undo", map[string]any{}, http.MethodPost); err != nil {
		t.Fatal(err)
	}
	missing, err := postJSON("api/task?id="+id, nil, http.MethodGet)
	if err != nil || missing["error"] != "task not found" {
		t.Fatalf("undo did not remove created task: %v", err)
	}
}
