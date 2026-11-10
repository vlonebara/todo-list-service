package db

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/vlonebara/todo-list/pkg/model"
)

// testStore создаёт отдельное хранилище для каждого теста.
func testStore(t *testing.T) TaskStore {
	t.Helper()
	database, err := OpenDatabase(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Error(err)
		}
	})
	return NewTaskStore(database)
}

// TestTaskStoreCanceledRequest проверяет отмену записи без изменения БД.
func TestTaskStoreCanceledRequest(t *testing.T) {
	store := testStore(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := store.Add(ctx, model.Task{Date: "20300101", Title: "test"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled operation, got %v", err)
	}
	tasks, err := store.List(t.Context(), model.TaskFilter{})
	if err != nil || len(tasks) != 0 {
		t.Fatalf("canceled request changed tasks: %v, %v", tasks, err)
	}
}
