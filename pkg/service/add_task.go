package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vlonebara/todo-list/pkg/model"
)

// AddTask проверяет и сохраняет задачу, возвращая её ID.
func (s TaskService) AddTask(ctx context.Context, date, title, comment, repeat string) (int64, error) {
	date, err := prepareTaskDate(date, title, repeat, time.Now())
	if err != nil {
		return 0, err
	}

	return s.store.Add(ctx, model.Task{Date: date, Title: title, Comment: comment, Repeat: repeat})
}

// prepareTaskDate проверяет поля задачи и исправляет просроченную дату.
func prepareTaskDate(date, title, repeat string, now time.Time) (string, error) {
	if strings.TrimSpace(title) == "" {
		return "", fmt.Errorf("%w: title is required", model.ErrInvalidTask)
	}

	today := now.Format(dateLayout)
	if date == "" {
		date = today
	}
	if _, err := time.Parse(dateLayout, date); err != nil {
		return "", fmt.Errorf("%w: invalid date: %v", model.ErrInvalidTask, err)
	}

	if repeat != "" {
		nextDate, err := NextDate(now, date, repeat)
		if err != nil {
			return "", fmt.Errorf("%w: invalid repeat rule: %v", model.ErrInvalidTask, err)
		}
		if date < today {
			date = nextDate
		}
	} else if date < today {
		date = today
	}
	return date, nil
}
