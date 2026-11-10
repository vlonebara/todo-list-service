package service

import (
	"context"
	"time"

	"github.com/vlonebara/todo-list/pkg/model"
)

// EditTask проверяет и сохраняет новые данные задачи.
func (s TaskService) EditTask(ctx context.Context, id, date, title, comment, repeat string) error {
	taskID, err := parseTaskID(id)
	if err != nil {
		return err
	}

	date, err = prepareTaskDate(date, title, repeat, time.Now())
	if err != nil {
		return err
	}

	return s.store.Update(ctx, model.Task{ID: taskID, Date: date, Title: title, Comment: comment, Repeat: repeat})
}
