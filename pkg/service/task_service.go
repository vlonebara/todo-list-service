package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

// TaskStore описывает операции хранилища, необходимые сервису задач.
type TaskStore interface {
	Add(ctx context.Context, task model.Task) (int64, error)
	Get(ctx context.Context, id int64) (model.Task, error)
	List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error)
	All(ctx context.Context) ([]model.Task, error)
	Ping(ctx context.Context) error
	Update(ctx context.Context, task model.Task) error
	Delete(ctx context.Context, id int64) error
	Complete(ctx context.Context, original model.Task, date string) error
	History(ctx context.Context, offset int) ([]model.TaskChange, error)
	Undo(ctx context.Context) error
}

// TaskService проверяет данные задач и обращается к общему хранилищу.
type TaskService struct {
	store TaskStore
}

// NewTaskService создаёт сервис с готовым хранилищем задач.
func NewTaskService(store TaskStore) TaskService {
	return TaskService{store: store}
}

// parseTaskID проверяет и преобразует ID задачи в число.
func parseTaskID(id string) (int64, error) {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if errutil.CheckError(err) || taskID < 1 {
		return 0, fmt.Errorf("%w: invalid task id", model.ErrInvalidTask)
	}
	return taskID, nil
}
