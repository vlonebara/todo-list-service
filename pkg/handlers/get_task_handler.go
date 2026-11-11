package handlers

import (
	"net/http"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// getTaskHandler возвращает данные задачи с указанным ID.
func (h Handler) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := h.tasks.GetTask(r.Context(), r.URL.Query().Get("id"))
	if errutil.CheckError(err) {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, taskToResponse(task))
}
