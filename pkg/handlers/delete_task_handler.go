package handlers

import "net/http"

// deleteTaskHandler удаляет задачу с указанным ID.
func (h Handler) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	err := h.tasks.DeleteTask(r.Context(), r.URL.Query().Get("id"))
	writeTaskResult(w, err)
}
