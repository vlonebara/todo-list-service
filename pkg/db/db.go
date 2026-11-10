package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var schema = []string{
	`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8),
		title TEXT,
		comment TEXT,
		repeat VARCHAR(128) CHECK (repeat IS NULL OR length(repeat) <= 128)
	)`,
	`CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date)`,
}

var historySchema = []string{
	`CREATE TABLE IF NOT EXISTS task_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		action TEXT NOT NULL,
		before_task TEXT NOT NULL,
		after_task TEXT NOT NULL,
		created_at TEXT NOT NULL,
		undone INTEGER NOT NULL DEFAULT 0
	)`,
}

// Create применяет новые версии схемы, сохраняя существующие задачи.
func Create(database *sql.DB) error {
	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("begin migration: %w", err)
	}
	defer tx.Rollback()
	var version int
	if err := tx.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}
	migrations := [][]string{schema, historySchema}
	if version > len(migrations) {
		return fmt.Errorf("database schema is newer than this application")
	}
	for ; version < len(migrations); version++ {
		for _, statement := range migrations[version] {
			if _, err := tx.Exec(statement); err != nil {
				return fmt.Errorf("migrate database: %w", err)
			}
		}
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", version)); err != nil {
		return fmt.Errorf("save schema version: %w", err)
	}
	return tx.Commit()
}
