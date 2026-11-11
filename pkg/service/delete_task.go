package service

import (
	"context"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// DeleteTask удаляет задачу с указанным ID.
func (s TaskService) DeleteTask(ctx context.Context, id string) error {
	taskID, err := parseTaskID(id)
	if errutil.CheckError(err) {
		return err
	}
	return s.store.Delete(ctx, taskID)
}
