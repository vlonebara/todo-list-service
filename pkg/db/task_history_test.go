package db

import (
	"errors"
	"sync"
	"testing"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

// TestTaskHistory проверяет историю и последовательную отмену действий.
func TestTaskHistory(t *testing.T) {
	store := testStore(t)
	ctx := t.Context()
	task := addTestTask(t, store, model.Task{Date: "20300101", Title: "original", Repeat: "d 1"})
	task.Title = "edited"
	if err := store.Update(ctx, task); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if err := store.Complete(ctx, task, "20300102"); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if err := store.Delete(ctx, task.ID); errutil.CheckError(err) {
		t.Fatal(err)
	}
	changes, err := store.History(ctx, 0)
	if errutil.CheckError(err) || len(changes) != 4 {
		t.Fatalf("unexpected history: %v, %v", changes, err)
	}
	if err := store.Undo(ctx); errutil.CheckError(err) {
		t.Fatal(err)
	}
	got, err := store.Get(ctx, task.ID)
	if errutil.CheckError(err) || got.Date != "20300102" {
		t.Fatalf("delete was not undone: %v, %v", got, err)
	}
	if err := store.Undo(ctx); errutil.CheckError(err) {
		t.Fatal(err)
	}
	got, err = store.Get(ctx, task.ID)
	if errutil.CheckError(err) || got != task {
		t.Fatalf("completion was not undone: %v, %v", got, err)
	}
	if err := store.Undo(ctx); errutil.CheckError(err) {
		t.Fatal(err)
	}
	got, err = store.Get(ctx, task.ID)
	if errutil.CheckError(err) || got.Title != "original" {
		t.Fatalf("edit was not undone: %v, %v", got, err)
	}
	if err := store.Undo(ctx); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if _, err := store.Get(ctx, task.ID); !errors.Is(err, model.ErrTaskNotFound) {
		t.Fatalf("add was not undone: %v", err)
	}
	if err := store.Undo(ctx); !errors.Is(err, model.ErrHistoryEmpty) {
		t.Fatalf("expected empty history: %v", err)
	}
}

// TestConcurrentCompletion запрещает два изменения на основе одного снимка.
func TestConcurrentCompletion(t *testing.T) {
	store := testStore(t)
	task := addTestTask(t, store, model.Task{Date: "20300101", Title: "test", Repeat: "d 1"})
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); results <- store.Complete(t.Context(), task, "20300102") }()
	}
	wg.Wait()
	close(results)
	var success, conflicts int
	for err := range results {
		switch {
		case !errutil.CheckError(err):
			success++
		case errors.Is(err, model.ErrConflict):
			conflicts++
		default:
			t.Fatal(err)
		}
	}
	if success != 1 || conflicts != 1 {
		t.Fatalf("success=%d, conflicts=%d", success, conflicts)
	}
	changes, err := store.History(t.Context(), 0)
	if errutil.CheckError(err) || len(changes) != 2 {
		t.Fatalf("duplicate completion in history: %v", err)
	}
}

// TestHistoryFailureRollback проверяет откат задачи при ошибке записи истории.
func TestHistoryFailureRollback(t *testing.T) {
	store := testStore(t)
	task := addTestTask(t, store, model.Task{Date: "20300101", Title: "original"})
	_, err := store.db.Exec(`CREATE TRIGGER reject_history BEFORE INSERT ON task_history BEGIN SELECT RAISE(ABORT, 'test failure'); END`)
	if errutil.CheckError(err) {
		t.Fatal(err)
	}
	changed := task
	changed.Title = "changed"
	if err := store.Update(t.Context(), changed); !errutil.CheckError(err) {
		t.Fatal("expected history failure")
	}
	got, err := store.Get(t.Context(), task.ID)
	if errutil.CheckError(err) || got != task {
		t.Fatalf("failed action changed task: %v, %v", got, err)
	}
}

// TestUndoDeletedTaskConflict не восстанавливает задачу поверх внешнего удаления.
func TestUndoDeletedTaskConflict(t *testing.T) {
	store := testStore(t)
	task := addTestTask(t, store, model.Task{Date: "20300101", Title: "test"})
	if _, err := store.db.Exec(`DELETE FROM scheduler WHERE id = ?`, task.ID); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if err := store.Undo(t.Context()); !errors.Is(err, model.ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}
