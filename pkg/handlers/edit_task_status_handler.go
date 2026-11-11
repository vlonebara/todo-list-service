package handlers

import "net/http"

// editTaskStatusHandler помечает задачу выполненной.
func (h Handler) editTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	err := h.tasks.EditTaskStatus(r.Context(), r.URL.Query().Get("id"))
	writeTaskResult(w, err)
}
