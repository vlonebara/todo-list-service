package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

// Undo восстанавливает состояние перед последним неотменённым действием.
func (s TaskStore) Undo(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if errutil.CheckError(err) {
		return fmt.Errorf("begin undo: %w", err)
	}
	defer tx.Rollback()
	var change model.TaskChange
	var before, after string
	err = tx.QueryRowContext(ctx,
		`SELECT id, before_task, after_task FROM task_history WHERE undone = 0 ORDER BY id DESC LIMIT 1`,
	).Scan(&change.ID, &before, &after)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrHistoryEmpty
	}
	if errutil.CheckError(err) {
		return fmt.Errorf("get last action: %w", err)
	}
	if err := decodeChange(&change, before, after); errutil.CheckError(err) {
		return err
	}
	if err := checkUndoState(ctx, tx, change); errutil.CheckError(err) {
		return err
	}
	if change.Before == nil {
		err = deleteTask(ctx, tx, change.After.ID)
	} else {
		task := change.Before
		_, err = tx.ExecContext(ctx,
			`INSERT INTO scheduler (id, date, title, comment, repeat) VALUES (?, ?, ?, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET date = excluded.date, title = excluded.title,
			 comment = excluded.comment, repeat = excluded.repeat`,
			task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	}
	if errutil.CheckError(err) {
		return fmt.Errorf("restore task: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE task_history SET undone = 1 WHERE id = ?`, change.ID); errutil.CheckError(err) {
		return fmt.Errorf("mark action undone: %w", err)
	}
	return tx.Commit()
}

// checkUndoState запрещает отмену поверх неизвестных изменений задачи.
func checkUndoState(ctx context.Context, tx *sql.Tx, change model.TaskChange) error {
	if change.Before == nil && change.After == nil {
		return model.ErrConflict
	}
	id := int64(0)
	if change.Before != nil {
		id = change.Before.ID
	}
	if change.After != nil {
		id = change.After.ID
	}
	current, err := scanTask(tx.QueryRowContext(ctx, taskByID, id))
	if change.After == nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil
		}
		if errutil.CheckError(err) {
			return err
		}
		return model.ErrConflict
	}
	if errutil.CheckError(err) {
		if errors.Is(err, model.ErrTaskNotFound) {
			return model.ErrConflict
		}
		return err
	}
	if current != *change.After {
		return model.ErrConflict
	}
	return nil
}
