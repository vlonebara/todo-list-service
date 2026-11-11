package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// openTestDatabase открывает тестовую БД и закрывает её после проверки.
func openTestDatabase(t *testing.T, file string) *sql.DB {
	t.Helper()
	database, err := OpenDatabase(file)
	if errutil.CheckError(err) {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := database.Close(); errutil.CheckError(err) {
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
	if errutil.CheckError(err) {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("unexpected database permissions: %v", info.Mode().Perm())
	}
}
