package handlers

import "net/http"

// getTaskHandler возвращает данные задачи с указанным ID.
func (h Handler) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := h.tasks.GetTask(r.URL.Query().Get("id"))
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, taskToResponse(task))
}
