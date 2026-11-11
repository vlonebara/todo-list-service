package handlers

import (
	"net/http"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// exportTasksHandler возвращает все задачи в загружаемом JSON-файле.
func (h Handler) exportTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.tasks.ExportTasks(r.Context())
	if errutil.CheckError(err) {
		writeTaskError(w, err)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="tasks.json"`)
	writeTasks(w, tasks)
}
