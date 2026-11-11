package service

import (
	"context"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

// EditTask проверяет и сохраняет новые данные задачи.
func (s TaskService) EditTask(ctx context.Context, id, date, title, comment, repeat string) error {
	if err := validateTaskText(title, comment); errutil.CheckError(err) {
		return err
	}
	taskID, err := parseTaskID(id)
	if errutil.CheckError(err) {
		return err
	}

	date, err = prepareTaskDate(date, title, repeat, time.Now())
	if errutil.CheckError(err) {
		return err
	}

	return s.store.Update(ctx, model.Task{ID: taskID, Date: date, Title: title, Comment: comment, Repeat: repeat})
}
