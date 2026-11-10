package service

import "context"

// DeleteTask удаляет задачу с указанным ID.
func (s TaskService) DeleteTask(ctx context.Context, id string) error {
	taskID, err := parseTaskID(id)
	if err != nil {
		return err
	}
	return s.store.Delete(ctx, taskID)
}
