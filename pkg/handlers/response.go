package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/vlonebara/todo-list/pkg/model"
)

type apiErrorResponse struct {
	Error string `json:"error"`
}

// writeTaskResult возвращает пустой JSON после изменения задачи или описывает ошибку.
func writeTaskResult(w http.ResponseWriter, err error) {
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}

// writeTaskError выбирает HTTP-код и скрывает внутренние ошибки сервера.
func writeTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidTask):
		writeAPIError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, model.ErrTaskNotFound):
		writeAPIError(w, http.StatusNotFound, "task not found")
	case errors.Is(err, model.ErrHistoryEmpty):
		writeAPIError(w, http.StatusNotFound, "history is empty")
	case errors.Is(err, model.ErrConflict):
		writeAPIError(w, http.StatusConflict, err.Error())
	default:
		log.Printf("task request: %v", err)
		writeAPIError(w, http.StatusInternalServerError, "internal server error")
	}
}

// writeAPIError возвращает клиенту JSON с описанием ошибки.
func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiErrorResponse{Error: message})
}

// writeJSON записывает значение в ответ в формате JSON.
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}
