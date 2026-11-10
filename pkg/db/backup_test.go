package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vlonebara/todo-list/pkg/model"
)

// TestBackupRestore проверяет восстановление задач и истории из копии.
func TestBackupRestore(t *testing.T) {
	store := testStore(t)
	id, err := store.Add(t.Context(), model.Task{Date: "20300101", Title: "saved task"})
	if err != nil {
		t.Fatal(err)
	}
	backup := filepath.Join(t.TempDir(), "backup.db")
	if err := Backup(t.Context(), store.db, backup); err != nil {
		t.Fatal(err)
	}
	if err := Backup(t.Context(), store.db, backup); err == nil {
		t.Fatal("backup overwrote an existing file")
	}
	info, err := os.Stat(backup)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected backup permissions: %v", err)
	}
	destination := filepath.Join(t.TempDir(), "restored.db")
	if err := Restore(t.Context(), backup, destination); err != nil {
		t.Fatal(err)
	}
	if err := Restore(t.Context(), backup, destination); err == nil {
		t.Fatal("restore overwrote an existing file")
	}
	database, err := OpenDatabase(destination)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	restored := NewTaskStore(database)
	task, err := restored.Get(t.Context(), id)
	if err != nil || task.Title != "saved task" {
		t.Fatalf("task was not restored: %v", err)
	}
	changes, err := restored.History(t.Context(), 0)
	if err != nil || len(changes) != 1 {
		t.Fatalf("history was not restored: %v", err)
	}
}

// TestRestoreRejectsInvalidFile проверяет отказ без создания целевой БД.
func TestRestoreRejectsInvalidFile(t *testing.T) {
	source := filepath.Join(t.TempDir(), "invalid.db")
	if err := os.WriteFile(source, []byte("not a database"), 0600); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "restored.db")
	if err := Restore(t.Context(), source, destination); err == nil {
		t.Fatal("invalid backup was accepted")
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatal("invalid restore created a destination")
	}
}
