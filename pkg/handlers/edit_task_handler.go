package handlers

import (
	"net/http"
)

// editTaskHandler обновляет задачу и возвращает ошибку, если обновление не удалось.
func (h Handler) editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request taskData
	if err := decodeJSONRequest(r, &request); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.tasks.EditTask(r.Context(), request.ID, request.Date, request.Title, request.Comment, request.Repeat); err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}
