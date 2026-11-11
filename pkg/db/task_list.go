package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

const taskFields = `SELECT id, date, title, comment, repeat FROM scheduler`

// List возвращает страницу задач с поиском без учёта регистра, включая кириллицу.
func (s TaskStore) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	query, args := tasksListQuery(filter)
	skip := filter.Offset
	if filter.Text == "" {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, model.PageSize, filter.Offset)
		skip = 0
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if errutil.CheckError(err) {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	defer rows.Close()
	return readTasks(rows, filter.Text, skip, model.PageSize)
}

// All возвращает все задачи для выгрузки пользователем.
func (s TaskStore) All(ctx context.Context) ([]model.Task, error) {
	rows, err := s.db.QueryContext(ctx, taskFields+` ORDER BY date, id`)
	if errutil.CheckError(err) {
		return nil, fmt.Errorf("export tasks: %w", err)
	}
	defer rows.Close()
	return readTasks(rows, "", 0, 0)
}

// readTasks читает строки и применяет текстовый поиск до разбиения на страницы.
func readTasks(rows *sql.Rows, search string, skip, limit int) ([]model.Task, error) {
	tasks := make([]model.Task, 0, model.PageSize)
	search = strings.ToLower(search)
	for rows.Next() {
		task, err := scanTask(rows)
		if errutil.CheckError(err) {
			return nil, err
		}
		// Встроенный SQLite LIKE не меняет регистр русских букв.
		if search != "" && !strings.Contains(strings.ToLower(task.Title), search) && !strings.Contains(strings.ToLower(task.Comment), search) {
			continue
		}
		if skip > 0 {
			skip--
			continue
		}
		tasks = append(tasks, task)
		if limit > 0 && len(tasks) == limit {
			break
		}
	}
	return tasks, rows.Err()
}

// tasksListQuery строит параметризованный запрос по календарным границам.
func tasksListQuery(filter model.TaskFilter) (string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0, 3)
	if filter.Date != "" {
		conditions = append(conditions, "date = ?")
		args = append(args, filter.Date)
	}
	if filter.FromDate != "" {
		conditions = append(conditions, "date >= ?")
		args = append(args, filter.FromDate)
	}
	if filter.BeforeDate != "" {
		conditions = append(conditions, "date < ?")
		args = append(args, filter.BeforeDate)
	}
	query := taskFields
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query + ` ORDER BY date, id`, args
}

// Ping проверяет доступность общего соединения с БД.
func (s TaskStore) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }
