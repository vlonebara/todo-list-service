package handlers

import "net/http"

type getTasksListResponse struct {
	Tasks []taskData `json:"tasks"`
}

// getTasksListHandler возвращает ближайшие задачи, отсортированные по дате.
func (h Handler) getTasksListHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.tasks.GetTasksList(r.Context(), r.URL.Query().Get("search"))
	if err != nil {
		writeTaskError(w, err)
		return
	}

	response := getTasksListResponse{Tasks: make([]taskData, 0, len(tasks))}
	for _, task := range tasks {
		response.Tasks = append(response.Tasks, taskToResponse(task))
	}
	writeJSON(w, http.StatusOK, response)
}
