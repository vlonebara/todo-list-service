package service

import (
	"context"
	"fmt"

	"github.com/vlonebara/todo-list/pkg/model"
)

// GetHistory возвращает страницу последних действий пользователя.
func (s TaskService) GetHistory(ctx context.Context, page int) ([]model.TaskChange, error) {
	offset, err := pageOffset(page)
	if err != nil {
		return nil, err
	}
	return s.store.History(ctx, offset)
}

// pageOffset проверяет номер страницы и вычисляет смещение выборки.
func pageOffset(page int) (int, error) {
	if page < 1 || page > 1000000 {
		return 0, fmt.Errorf("%w: page must be between 1 and 1000000", model.ErrInvalidTask)
	}
	return (page - 1) * model.PageSize, nil
}

// UndoLast отменяет последнее действие, если задача ещё не изменена извне.
func (s TaskService) UndoLast(ctx context.Context) error {
	return s.store.Undo(ctx)
}
