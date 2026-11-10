package service

import (
	"context"
	"github.com/vlonebara/todo-list/pkg/model"
)

// ExportTasks получает полный список задач для сохранения в файл.
func (s TaskService) ExportTasks(ctx context.Context) ([]model.Task, error) {
	return s.store.All(ctx)
}

// Health проверяет готовность хранилища принимать запросы.
func (s TaskService) Health(ctx context.Context) error {
	return s.store.Ping(ctx)
}
