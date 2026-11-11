package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

const taskByID = taskFields + ` WHERE id = ?`

// rowScanner позволяет читать задачу из одной строки или результата выборки.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanTask читает одну задачу и распознаёт отсутствие записи.
func scanTask(row rowScanner) (model.Task, error) {
	var task model.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return task, model.ErrTaskNotFound
	}
	if errutil.CheckError(err) {
		return task, fmt.Errorf("read task: %w", err)
	}
	return task, nil
}

// beginChange открывает транзакцию и читает прежнее состояние задачи.
func (s TaskStore) beginChange(ctx context.Context, id int64) (*sql.Tx, *model.Task, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if errutil.CheckError(err) {
		return nil, nil, fmt.Errorf("begin task change: %w", err)
	}
	if id == 0 {
		return tx, nil, nil
	}
	task, err := scanTask(tx.QueryRowContext(ctx, taskByID, id))
	if errutil.CheckError(err) {
		_ = tx.Rollback()
		return nil, nil, err
	}
	return tx, &task, nil
}

// finishChange сохраняет историю и подтверждает изменение одной транзакцией.
func finishChange(ctx context.Context, tx *sql.Tx, action string, before, after *model.Task) error {
	oldJSON, err := json.Marshal(before)
	if errutil.CheckError(err) {
		return fmt.Errorf("encode previous task: %w", err)
	}
	newJSON, err := json.Marshal(after)
	if errutil.CheckError(err) {
		return fmt.Errorf("encode changed task: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO task_history (action, before_task, after_task, created_at) VALUES (?, ?, ?, ?)`,
		action, string(oldJSON), string(newJSON), time.Now().UTC().Format(time.RFC3339Nano))
	if errutil.CheckError(err) {
		return fmt.Errorf("save task history: %w", err)
	}
	return tx.Commit()
}

// History возвращает страницу последних действий вместе с отменёнными.
func (s TaskStore) History(ctx context.Context, offset int) ([]model.TaskChange, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, action, before_task, after_task, created_at, undone FROM task_history ORDER BY id DESC LIMIT ? OFFSET ?`, model.PageSize, offset)
	if errutil.CheckError(err) {
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()
	changes := make([]model.TaskChange, 0, model.PageSize)
	for rows.Next() {
		var change model.TaskChange
		var before, after string
		if err := rows.Scan(&change.ID, &change.Action, &before, &after, &change.CreatedAt, &change.Undone); errutil.CheckError(err) {
			return nil, fmt.Errorf("read history: %w", err)
		}
		if err := decodeChange(&change, before, after); errutil.CheckError(err) {
			return nil, err
		}
		changes = append(changes, change)
	}
	return changes, rows.Err()
}

// decodeChange восстанавливает снимки задачи из записи истории.
func decodeChange(change *model.TaskChange, before, after string) error {
	if err := json.Unmarshal([]byte(before), &change.Before); errutil.CheckError(err) {
		return fmt.Errorf("decode previous task: %w", err)
	}
	if err := json.Unmarshal([]byte(after), &change.After); errutil.CheckError(err) {
		return fmt.Errorf("decode changed task: %w", err)
	}
	return nil
}
