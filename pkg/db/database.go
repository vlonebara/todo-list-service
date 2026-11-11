package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// OpenDatabase создаёт общий пул соединений и подготавливает БД.
func OpenDatabase(file string) (*sql.DB, error) {
	if err := prepareDatabaseFile(file); err != nil {
		return nil, err
	}
	database, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Одно соединение упорядочивает записи SQLite и переиспользуется запросами.
	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)

	if err := configureDatabase(database); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}

// configureDatabase проверяет соединение и подготавливает настройки и схему.
func configureDatabase(database *sql.DB) error {
	if err := database.Ping(); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	if _, err := database.Exec(`PRAGMA busy_timeout = 5000`); err != nil {
		return fmt.Errorf("configure database: %w", err)
	}
	return Create(database)
}

// prepareDatabaseFile создаёт новый файл с доступом только для владельца.
func prepareDatabaseFile(path string) error {
	if path == ":memory:" || strings.HasPrefix(path, "file:") {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("create database file: %w", err)
	}
	return file.Close()
}
