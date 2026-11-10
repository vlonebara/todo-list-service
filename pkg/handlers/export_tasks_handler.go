package handlers

import "net/http"

// exportTasksHandler возвращает все задачи в загружаемом JSON-файле.
func (h Handler) exportTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.tasks.ExportTasks(r.Context())
	if err != nil {
		writeTaskError(w, err)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="tasks.json"`)
	writeTasks(w, tasks)
}
