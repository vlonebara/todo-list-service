package handlers

import (
	"net/http"
)

// editTaskHandler обновляет задачу и возвращает ошибку, если обновление не удалось.
func (h Handler) editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request taskData
	if !readJSONRequest(w, r, &request) {
		return
	}

	err := h.tasks.EditTask(r.Context(), request.ID, request.Date, request.Title, request.Comment, request.Repeat)
	writeTaskResult(w, err)
}
