package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

// openTestDatabase открывает тестовую БД и закрывает её после проверки.
func openTestDatabase(t *testing.T, file string) *sql.DB {
	t.Helper()
	database, err := OpenDatabase(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Error(err)
		}
	})
	return database
}

// TestDatabaseFilePermissions проверяет доступ к новой БД только для владельца.
func TestDatabaseFilePermissions(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.db")
	openTestDatabase(t, file)
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected database permissions: %v", info.Mode().Perm())
	}
}
