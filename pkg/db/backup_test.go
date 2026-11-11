package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

// TestBackupRestore проверяет восстановление задач и истории из копии.
func TestBackupRestore(t *testing.T) {
	store := testStore(t)
	saved := addTestTask(t, store, model.Task{Date: "20300101", Title: "saved task"})
	backup := filepath.Join(t.TempDir(), "backup.db")
	if err := Backup(t.Context(), store.db, backup); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if err := Backup(t.Context(), store.db, backup); !errutil.CheckError(err) {
		t.Fatal("backup overwrote an existing file")
	}
	info, err := os.Stat(backup)
	if errutil.CheckError(err) || info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected backup permissions: %v", err)
	}
	destination := filepath.Join(t.TempDir(), "restored.db")
	if err := Restore(t.Context(), backup, destination); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if err := Restore(t.Context(), backup, destination); !errutil.CheckError(err) {
		t.Fatal("restore overwrote an existing file")
	}
	database := openTestDatabase(t, destination)
	restored := NewTaskStore(database)
	task, err := restored.Get(t.Context(), saved.ID)
	if errutil.CheckError(err) || task.Title != "saved task" {
		t.Fatalf("task was not restored: %v", err)
	}
	changes, err := restored.History(t.Context(), 0)
	if errutil.CheckError(err) || len(changes) != 1 {
		t.Fatalf("history was not restored: %v", err)
	}
}

// TestRestoreRejectsInvalidFile проверяет отказ без создания целевой БД.
func TestRestoreRejectsInvalidFile(t *testing.T) {
	source := filepath.Join(t.TempDir(), "invalid.db")
	if err := os.WriteFile(source, []byte("not a database"), 0600); errutil.CheckError(err) {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "restored.db")
	if err := Restore(t.Context(), source, destination); !errutil.CheckError(err) {
		t.Fatal("invalid backup was accepted")
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatal("invalid restore created a destination")
	}
}
