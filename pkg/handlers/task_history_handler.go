package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vlonebara/todo-list/pkg/model"
)

type historyData struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Before    *taskData `json:"before"`
	After     *taskData `json:"after"`
	CreatedAt string    `json:"created_at"`
	Undone    bool      `json:"undone"`
}

// getHistoryHandler возвращает историю создания, изменения и выполнения задач.
func (h Handler) getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	page, err := requestPage(r)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	changes, err := h.tasks.GetHistory(r.Context(), page)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	items := make([]historyData, 0, len(changes))
	for _, change := range changes {
		item := historyData{ID: strconv.FormatInt(change.ID, 10), Action: change.Action,
			CreatedAt: change.CreatedAt, Undone: change.Undone}
		if change.Before != nil {
			task := taskToResponse(*change.Before)
			item.Before = &task
		}
		if change.After != nil {
			task := taskToResponse(*change.After)
			item.After = &task
		}
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, map[string]any{"history": items})
}

// undoTaskHandler отменяет последнее действие над задачами.
func (h Handler) undoTaskHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.tasks.UndoLast(r.Context()); err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}

// requestPage читает номер страницы и использует первую страницу по умолчанию.
func requestPage(r *http.Request) (int, error) {
	value := r.URL.Query().Get("page")
	if value == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(value)
	if err != nil || page < 1 || page > 1000000 {
		return 0, fmt.Errorf("%w: page must be between 1 and 1000000", model.ErrInvalidTask)
	}
	return page, nil
}
