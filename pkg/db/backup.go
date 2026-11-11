package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// Backup сохраняет согласованную копию БД, не заменяя существующий файл.
func Backup(ctx context.Context, database *sql.DB, destination string) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create backup file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close backup file: %w", err)
	}
	if _, err := database.ExecContext(ctx, `VACUUM main INTO ?`, destination); err != nil {
		return fmt.Errorf("backup database: %w", err)
	}
	return nil
}

// Restore проверяет копию и восстанавливает её в новый файл БД.
func Restore(ctx context.Context, source, destination string) error {
	if info, err := os.Stat(source); err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("backup source must be an existing regular file")
	}
	database, err := sql.Open("sqlite", source)
	if err != nil {
		return fmt.Errorf("open backup: %w", err)
	}
	defer database.Close()
	var result string
	if err := database.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&result); err != nil {
		return fmt.Errorf("check backup: %w", err)
	}
	if result != "ok" {
		return fmt.Errorf("backup database is corrupted")
	}
	if err := validateBackup(ctx, database); err != nil {
		return err
	}
	return Backup(ctx, database, destination)
}

// validateBackup проверяет, что копия содержит таблицу задач приложения.
func validateBackup(ctx context.Context, database *sql.DB) error {
	rows, err := database.QueryContext(ctx, taskFields+` LIMIT 0`)
	if err != nil {
		return fmt.Errorf("backup does not contain the task schema: %w", err)
	}
	return rows.Close()
}
