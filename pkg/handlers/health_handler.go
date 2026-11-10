package handlers

import (
	"context"
	"net/http"
	"time"
)

// healthHandler сообщает о готовности сервера без раскрытия ошибок БД.
func (h Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.tasks.Health(ctx); err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, "service is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
