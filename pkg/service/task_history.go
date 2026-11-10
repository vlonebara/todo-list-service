package service

import (
	"context"
	"fmt"

	"github.com/vlonebara/todo-list/pkg/model"
)

// GetHistory возвращает страницу последних действий пользователя.
func (s TaskService) GetHistory(ctx context.Context, page int) ([]model.TaskChange, error) {
	if page < 1 || page > 1000000 {
		return nil, fmt.Errorf("%w: page must be between 1 and 1000000", model.ErrInvalidTask)
	}
	return s.store.History(ctx, (page-1)*15)
}

// UndoLast отменяет последнее действие, если задача ещё не изменена извне.
func (s TaskService) UndoLast(ctx context.Context) error {
	return s.store.Undo(ctx)
}
