package db

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDatabaseFilePermissions проверяет доступ к новой БД только для владельца.
func TestDatabaseFilePermissions(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.db")
	database, err := OpenDatabase(file)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected database permissions: %v", info.Mode().Perm())
	}
}
