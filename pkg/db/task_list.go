package db

import (
	"context"
	"fmt"

	"github.com/vlonebara/todo-list/pkg/model"
)

// List возвращает до 15 задач по дате или тексту поиска.
func (s TaskStore) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	query, args := tasksListQuery(filter)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0, 15)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("read task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}
	return tasks, nil
}

// tasksListQuery строит запрос для списка, даты или текстового поиска.
func tasksListQuery(filter model.TaskFilter) (string, []any) {
	const fields = `SELECT id, date, title, comment, repeat FROM scheduler`
	const orderAndLimit = ` ORDER BY date, id LIMIT 15`
	if filter.Date != "" {
		return fields + ` WHERE date = ?` + orderAndLimit, []any{filter.Date}
	}
	if filter.Text != "" {
		pattern := "%" + filter.Text + "%"
		return fields + ` WHERE title LIKE ? OR comment LIKE ?` + orderAndLimit, []any{pattern, pattern}
	}
	return fields + orderAndLimit, nil
}
