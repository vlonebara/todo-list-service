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
	GetTasksList(ctx context.Context, search string) ([]model.Task, error)
	EditTask(ctx context.Context, id, date, title, comment, repeat string) error
	DeleteTask(ctx context.Context, id string) error
	EditTaskStatus(ctx context.Context, id string) error
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
}

// New создаёт HTTP-маршруты API и файловый сервер веб-интерфейса.
func New(tasks TaskService, authentication AuthService) *http.ServeMux {
	handlers := Handler{tasks: tasks, authentication: authentication}
	mux := http.NewServeMux()
	handlers.registerRoutes(mux)
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	return mux
}
