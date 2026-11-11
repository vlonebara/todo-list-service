package handlers

import "net/http"

// registerRoutes связывает маршруты API с обработчиками и middleware.
func (h Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/session", methods(map[string]http.HandlerFunc{
		http.MethodGet: h.sessionHandler,
	}))
	mux.HandleFunc("/api/health", methods(map[string]http.HandlerFunc{
		http.MethodGet: h.healthHandler,
	}))
	mux.HandleFunc("/api/tasks/export", methods(map[string]http.HandlerFunc{
		http.MethodGet: auth(h.authentication, h.exportTasksHandler),
	}))
	mux.HandleFunc("/api/history", methods(map[string]http.HandlerFunc{
		http.MethodGet: auth(h.authentication, h.getHistoryHandler),
	}))
	mux.HandleFunc("/api/task/undo", methods(map[string]http.HandlerFunc{
		http.MethodPost: auth(h.authentication, h.undoTaskHandler),
	}))
	mux.HandleFunc("/api/nextdate", methods(map[string]http.HandlerFunc{
		http.MethodGet: nextDayHandler,
	}))
	mux.HandleFunc("/api/signin", methods(map[string]http.HandlerFunc{
		http.MethodPost: limitSignin(h.authHandler),
	}))
	mux.HandleFunc("/api/signout", methods(map[string]http.HandlerFunc{
		http.MethodPost: h.signoutHandler,
	}))
	mux.HandleFunc("/api/tasks", methods(map[string]http.HandlerFunc{
		http.MethodGet: auth(h.authentication, h.getTasksListHandler),
	}))
	mux.HandleFunc("/api/task", methods(map[string]http.HandlerFunc{
		http.MethodGet:    auth(h.authentication, h.getTaskHandler),
		http.MethodPost:   auth(h.authentication, h.addTaskHandler),
		http.MethodPut:    auth(h.authentication, h.editTaskHandler),
		http.MethodDelete: auth(h.authentication, h.deleteTaskHandler),
	}))
	mux.HandleFunc("/api/task/done", methods(map[string]http.HandlerFunc{
		http.MethodPost: auth(h.authentication, h.editTaskStatusHandler),
	}))
}
