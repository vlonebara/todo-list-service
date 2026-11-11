package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// TestCreatePreservesLegacyDatabase проверяет обновление старой БД без потери задач.
func TestCreatePreservesLegacyDatabase(t *testing.T) {
	file := filepath.Join(t.TempDir(), "legacy.db")
	legacy, err := sql.Open("sqlite", file)
	if errutil.CheckError(err) {
		t.Fatal(err)
	}
	for _, statement := range schema {
		if _, err := legacy.Exec(statement); errutil.CheckError(err) {
			t.Fatal(err)
		}
	}
	if _, err := legacy.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES ('20300101', 'legacy', '', '')`); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if err := legacy.Close(); errutil.CheckError(err) {
		t.Fatal(err)
	}
	database := openTestDatabase(t, file)
	if err := Create(database); errutil.CheckError(err) {
		t.Fatal(err)
	}
	task, err := NewTaskStore(database).Get(t.Context(), 1)
	if errutil.CheckError(err) || task.Title != "legacy" {
		t.Fatalf("legacy task was changed: %v, %v", task, err)
	}
	var version int
	if err := database.QueryRow(`PRAGMA user_version`).Scan(&version); errutil.CheckError(err) {
		t.Fatal(err)
	}
	if version != 2 {
		t.Fatalf("unexpected schema version: %d", version)
	}
}
