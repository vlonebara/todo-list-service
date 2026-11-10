package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// TestCreatePreservesLegacyDatabase проверяет обновление старой БД без потери задач.
func TestCreatePreservesLegacyDatabase(t *testing.T) {
	file := filepath.Join(t.TempDir(), "legacy.db")
	legacy, err := sql.Open("sqlite", file)
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range schema {
		if _, err := legacy.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := legacy.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES ('20300101', 'legacy', '', '')`); err != nil {
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}
	database, err := OpenDatabase(file)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := Create(database); err != nil {
		t.Fatal(err)
	}
	task, err := NewTaskStore(database).Get(t.Context(), 1)
	if err != nil || task.Title != "legacy" {
		t.Fatalf("legacy task was changed: %v, %v", task, err)
	}
	var version int
	if err := database.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version != 2 {
		t.Fatalf("unexpected schema version: %d", version)
	}
}
