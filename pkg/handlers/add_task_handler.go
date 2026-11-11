package handlers

import (
	"net/http"
	"strconv"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// addTaskHandler обрабатывает POST-запрос и возвращает ID задачи или ошибку.
func (h Handler) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request taskData
	if !readJSONRequest(w, r, &request) {
		return
	}

	id, err := h.tasks.AddTask(r.Context(), request.Date, request.Title, request.Comment, request.Repeat)
	if errutil.CheckError(err) {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
