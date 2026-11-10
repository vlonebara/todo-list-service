package service

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/vlonebara/todo-list/pkg/model"
)

const searchDateLayout = "02.01.2006"

// GetTasksList возвращает до 15 задач по дате или тексту поиска.
func (s TaskService) GetTasksList(ctx context.Context, search string) ([]model.Task, error) {
	if utf8.RuneCountInString(search) > 256 {
		return nil, fmt.Errorf("%w: search must not exceed 256 characters", model.ErrInvalidTask)
	}
	return s.store.List(ctx, taskFilter(search))
}

// taskFilter распознаёт дату поиска и приводит её к формату задачи.
func taskFilter(search string) model.TaskFilter {
	date, err := time.Parse(searchDateLayout, search)
	if err == nil && date.Format(searchDateLayout) == search {
		return model.TaskFilter{Date: date.Format(dateLayout)}
	}
	return model.TaskFilter{Text: search}
}
