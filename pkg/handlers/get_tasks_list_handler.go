package handlers

import (
	"net/http"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/model"
)

type getTasksListResponse struct {
	Tasks []taskData `json:"tasks"`
}

// getTasksListHandler возвращает ближайшие задачи, отсортированные по дате.
func (h Handler) getTasksListHandler(w http.ResponseWriter, r *http.Request) {
	page, err := requestPage(r)
	if errutil.CheckError(err) {
		writeTaskError(w, err)
		return
	}
	tasks, err := h.tasks.GetTasksList(r.Context(), r.URL.Query().Get("search"), r.URL.Query().Get("period"), page)
	if errutil.CheckError(err) {
		writeTaskError(w, err)
		return
	}

	writeTasks(w, tasks)
}

// writeTasks возвращает общий JSON-формат списка и экспорта задач.
func writeTasks(w http.ResponseWriter, tasks []model.Task) {
	response := getTasksListResponse{Tasks: make([]taskData, 0, len(tasks))}
	for _, task := range tasks {
		response.Tasks = append(response.Tasks, taskToResponse(task))
	}
	writeJSON(w, http.StatusOK, response)
}
