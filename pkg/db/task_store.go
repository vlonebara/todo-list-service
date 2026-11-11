package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vlonebara/todo-list/pkg/model"
)

// TaskStore выполняет запросы к задачам через общий пул соединений.
type TaskStore struct {
	db *sql.DB
}

// NewTaskStore создаёт хранилище с уже открытым пулом соединений.
func NewTaskStore(database *sql.DB) TaskStore {
	return TaskStore{db: database}
}

// Add сохраняет задачу и возвращает её ID.
func (s TaskStore) Add(ctx context.Context, task model.Task) (int64, error) {
	tx, _, err := s.beginChange(ctx, 0)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx,
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	task.ID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get task id: %w", err)
	}
	if err := finishChange(ctx, tx, "add", nil, &task); err != nil {
		return 0, err
	}
	return task.ID, nil
}

// Get возвращает задачу с указанным ID.
func (s TaskStore) Get(ctx context.Context, id int64) (model.Task, error) {
	return scanTask(s.db.QueryRowContext(ctx, taskByID, id))
}

// Update сохраняет новые значения полей задачи.
func (s TaskStore) Update(ctx context.Context, task model.Task) error {
	tx, before, err := s.beginChange(ctx, task.ID)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx,
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return finishChange(ctx, tx, "edit", before, &task)
}

// Delete удаляет задачу с указанным ID.
func (s TaskStore) Delete(ctx context.Context, id int64) error {
	tx, before, err := s.beginChange(ctx, id)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM scheduler WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return finishChange(ctx, tx, "delete", before, nil)
}

// Complete меняет только ту задачу, на основе которой была рассчитана дата.
func (s TaskStore) Complete(ctx context.Context, original model.Task, date string) error {
	tx, before, err := s.beginChange(ctx, original.ID)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return model.ErrConflict
		}
		return err
	}
	defer tx.Rollback()
	if *before != original {
		return model.ErrConflict
	}
	var after *model.Task
	if date == "" {
		_, err = tx.ExecContext(ctx, `DELETE FROM scheduler WHERE id = ?`, original.ID)
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE scheduler SET date = ? WHERE id = ?`, date, original.ID)
		original.Date = date
		after = &original
	}
	if err != nil {
		return fmt.Errorf("complete task: %w", err)
	}
	return finishChange(ctx, tx, "complete", before, after)
}
