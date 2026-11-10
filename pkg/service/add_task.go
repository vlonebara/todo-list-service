package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vlonebara/todo-list/pkg/model"
)

// AddTask проверяет и сохраняет задачу, возвращая её ID.
func (s TaskService) AddTask(ctx context.Context, date, title, comment, repeat string) (int64, error) {
	if err := validateTaskText(title, comment); err != nil {
		return 0, err
	}
	date, err := prepareTaskDate(date, title, repeat, time.Now())
	if err != nil {
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
	parsed, err := time.Parse(dateLayout, date)
	if err != nil {
		return "", fmt.Errorf("%w: invalid date: %v", model.ErrInvalidTask, err)
	}
	if parsed.Year() < 1 {
		return "", fmt.Errorf("%w: date year must be between 1 and 9999", model.ErrInvalidTask)
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
