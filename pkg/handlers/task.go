package handlers

import (
	"strconv"

	"github.com/vlonebara/todo-list/pkg/model"
)

// taskData задаёт общий JSON-формат задачи для запросов и ответов.
type taskData struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// taskSnapshot преобразует снимок истории, сохраняя отсутствие задачи как null.
func taskSnapshot(task *model.Task) *taskData {
	if task == nil {
		return nil
	}
	response := taskToResponse(*task)
	return &response
}

// taskToResponse преобразует задачу сервиса в JSON-формат API.
func taskToResponse(task model.Task) taskData {
	return taskData{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
