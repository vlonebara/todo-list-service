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
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get task id: %w", err)
	}
	return id, nil
}

// Get возвращает задачу с указанным ID.
func (s TaskStore) Get(ctx context.Context, id int64) (model.Task, error) {
	var task model.Task
	err := s.db.QueryRowContext(ctx,
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, model.ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, fmt.Errorf("get task: %w", err)
	}
	return task, nil
}

// Update сохраняет новые значения полей задачи.
func (s TaskStore) Update(ctx context.Context, task model.Task) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return checkTaskChanged(result)
}

// Delete удаляет задачу с указанным ID.
func (s TaskStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return checkTaskChanged(result)
}

// SetDate записывает дату следующего выполнения задачи.
func (s TaskStore) SetDate(ctx context.Context, id int64, date string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE scheduler SET date = ? WHERE id = ?`, date, id)
	if err != nil {
		return fmt.Errorf("update task date: %w", err)
	}
	return checkTaskChanged(result)
}

// checkTaskChanged проверяет, что запрос затронул существующую задачу.
func checkTaskChanged(result sql.Result) error {
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check changed task: %w", err)
	}
	if changed == 0 {
		return model.ErrTaskNotFound
	}
	return nil
}
