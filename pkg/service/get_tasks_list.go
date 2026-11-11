package service

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

const searchDateLayout = "02.01.2006"

// GetTasksList возвращает до 15 задач по дате или тексту поиска.
func (s TaskService) GetTasksList(ctx context.Context, search, period string, page int) ([]model.Task, error) {
	if utf8.RuneCountInString(search) > 256 {
		return nil, fmt.Errorf("%w: search must not exceed 256 characters", model.ErrInvalidTask)
	}
	offset, err := pageOffset(page)
	if errutil.CheckError(err) {
		return nil, err
	}
	filter := taskFilter(search)
	filter.Offset = offset
	now := time.Now()
	switch period {
	case "", "all":
	case "today":
		filter.FromDate = now.Format(dateLayout)
		filter.BeforeDate = now.AddDate(0, 0, 1).Format(dateLayout)
	case "overdue":
		filter.BeforeDate = now.Format(dateLayout)
	case "upcoming":
		filter.FromDate = now.AddDate(0, 0, 1).Format(dateLayout)
	default:
		return nil, fmt.Errorf("%w: unknown task period", model.ErrInvalidTask)
	}
	return s.store.List(ctx, filter)
}

// taskFilter распознаёт дату поиска и приводит её к формату задачи.
func taskFilter(search string) model.TaskFilter {
	date, err := time.Parse(searchDateLayout, search)
	if !errutil.CheckError(err) && date.Format(searchDateLayout) == search {
		return model.TaskFilter{Date: date.Format(dateLayout)}
	}
	return model.TaskFilter{Text: search}
}
