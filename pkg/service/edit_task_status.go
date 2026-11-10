package service

import (
	"context"
	"fmt"
	"time"
)

// EditTaskStatus помечает задачу выполненной и переносит её или удаляет.
func (s TaskService) EditTaskStatus(ctx context.Context, id string) error {
	taskID, err := parseTaskID(id)
	if err != nil {
		return err
	}

	task, err := s.store.Get(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Repeat == "" {
		return s.store.Delete(ctx, taskID)
	}

	date, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("calculate next task date: %w", err)
	}
	return s.store.SetDate(ctx, taskID, date)
}
