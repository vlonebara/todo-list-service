package handlers

import (
	"context"
	"net/http"

	"github.com/vlonebara/todo-list/pkg/model"
)

// TaskService описывает операции задач, необходимые HTTP-обработчикам.
type TaskService interface {
	AddTask(ctx context.Context, date, title, comment, repeat string) (int64, error)
	GetTask(ctx context.Context, id string) (model.Task, error)
	GetTasksList(ctx context.Context, search, period string, page int) ([]model.Task, error)
	ExportTasks(ctx context.Context) ([]model.Task, error)
	Health(ctx context.Context) error
	EditTask(ctx context.Context, id, date, title, comment, repeat string) error
	DeleteTask(ctx context.Context, id string) error
	EditTaskStatus(ctx context.Context, id string) error
	GetHistory(ctx context.Context, page int) ([]model.TaskChange, error)
	UndoLast(ctx context.Context) error
}

// AuthService описывает проверку пароля и токена для HTTP-запросов.
type AuthService interface {
	Enabled() bool
	Authenticate(password string) (string, error)
	ValidateToken(token string) bool
}

// Handler передаёт HTTP-запросы сервисам через их интерфейсы.
type Handler struct {
	tasks          TaskService
	authentication AuthService
	secureCookies  bool
}

// New создаёт HTTP-маршруты API и файловый сервер веб-интерфейса.
func New(tasks TaskService, authentication AuthService, secureCookies bool, webDir string) http.Handler {
	handlers := Handler{tasks: tasks, authentication: authentication, secureCookies: secureCookies}
	mux := http.NewServeMux()
	handlers.registerRoutes(mux)
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	protection := http.NewCrossOriginProtection()
	protection.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusForbidden, "cross-origin request is not allowed")
	}))
	return securityHeaders(protection.Handler(mux))
}
