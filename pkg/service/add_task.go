package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

// AddTask проверяет и сохраняет задачу, возвращая её ID.
func (s TaskService) AddTask(ctx context.Context, date, title, comment, repeat string) (int64, error) {
	if err := validateTaskText(title, comment); errutil.CheckError(err) {
		return 0, err
	}
	date, err := prepareTaskDate(date, title, repeat, time.Now())
	if errutil.CheckError(err) {
		return 0, err
	}

	return s.store.Add(ctx, model.Task{Date: date, Title: title, Comment: comment, Repeat: repeat})
}

// validateTaskText ограничивает длину текстовых полей задачи.
func validateTaskText(title, comment string) error {
	if utf8.RuneCountInString(title) > 200 {
		return fmt.Errorf("%w: title must not exceed 200 characters", model.ErrInvalidTask)
	}
	if utf8.RuneCountInString(comment) > 5000 {
		return fmt.Errorf("%w: comment must not exceed 5000 characters", model.ErrInvalidTask)
	}
	return nil
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
	if _, err := ParseDate(date); errutil.CheckError(err) {
		return "", fmt.Errorf("%w: %v", model.ErrInvalidTask, err)
	}

	if repeat != "" {
		nextDate, err := NextDate(now, date, repeat)
		if errutil.CheckError(err) {
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
