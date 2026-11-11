package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// healthHandler сообщает о готовности сервера без раскрытия ошибок БД.
func (h Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.tasks.Health(ctx); errutil.CheckError(err) {
		writeAPIError(w, http.StatusServiceUnavailable, "service is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
