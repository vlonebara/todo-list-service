package service

import (
	"context"
	"github.com/vlonebara/todo-list/pkg/model"
)

// GetTask возвращает задачу с указанным ID.
func (s TaskService) GetTask(ctx context.Context, id string) (model.Task, error) {
	taskID, err := parseTaskID(id)
	if err != nil {
		return model.Task{}, err
	}
	return s.store.Get(ctx, taskID)
}
