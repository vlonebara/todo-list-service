package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/vlonebara/todo-list/pkg/model"
)

const taskByID = `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

// scanTask читает одну задачу и распознаёт отсутствие записи.
func scanTask(row *sql.Row) (model.Task, error) {
	var task model.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return task, model.ErrTaskNotFound
	}
	if err != nil {
		return task, fmt.Errorf("read task: %w", err)
	}
	return task, nil
}

// beginChange открывает транзакцию и читает прежнее состояние задачи.
func (s TaskStore) beginChange(ctx context.Context, id int64) (*sql.Tx, *model.Task, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("begin task change: %w", err)
	}
	if id == 0 {
		return tx, nil, nil
	}
	task, err := scanTask(tx.QueryRowContext(ctx, taskByID, id))
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}
	return tx, &task, nil
}

// finishChange сохраняет историю и подтверждает изменение одной транзакцией.
func finishChange(ctx context.Context, tx *sql.Tx, action string, before, after *model.Task) error {
	oldJSON, err := json.Marshal(before)
	if err != nil {
		return fmt.Errorf("encode previous task: %w", err)
	}
	newJSON, err := json.Marshal(after)
	if err != nil {
		return fmt.Errorf("encode changed task: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO task_history (action, before_task, after_task, created_at) VALUES (?, ?, ?, ?)`,
		action, string(oldJSON), string(newJSON), time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("save task history: %w", err)
	}
	return tx.Commit()
}

// History возвращает страницу последних действий вместе с отменёнными.
func (s TaskStore) History(ctx context.Context, offset int) ([]model.TaskChange, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, action, before_task, after_task, created_at, undone FROM task_history ORDER BY id DESC LIMIT ? OFFSET ?`, model.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()
	changes := make([]model.TaskChange, 0, model.PageSize)
	for rows.Next() {
		var change model.TaskChange
		var before, after string
		if err := rows.Scan(&change.ID, &change.Action, &before, &after, &change.CreatedAt, &change.Undone); err != nil {
			return nil, fmt.Errorf("read history: %w", err)
		}
		if err := decodeChange(&change, before, after); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	return changes, rows.Err()
}

// decodeChange восстанавливает снимки задачи из записи истории.
func decodeChange(change *model.TaskChange, before, after string) error {
	if err := json.Unmarshal([]byte(before), &change.Before); err != nil {
		return fmt.Errorf("decode previous task: %w", err)
	}
	if err := json.Unmarshal([]byte(after), &change.After); err != nil {
		return fmt.Errorf("decode changed task: %w", err)
	}
	return nil
}
