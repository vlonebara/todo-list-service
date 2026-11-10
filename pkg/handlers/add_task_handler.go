package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

// addTaskHandler обрабатывает POST-запрос и возвращает ID задачи или ошибку.
func (h Handler) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request taskData
	if err := decodeJSONRequest(r, &request); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.tasks.AddTask(request.Date, request.Title, request.Comment, request.Repeat)
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

// decodeJSONRequest читает один JSON-объект из тела запроса.
func decodeJSONRequest(r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(value); err != nil {
		return errors.New("invalid json")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("invalid json")
	}
	return nil
}
