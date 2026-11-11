package model

import "errors"

// ErrInvalidTask сообщает об ошибке в данных задачи.
var ErrInvalidTask = errors.New("invalid task")

// ErrTaskNotFound сообщает, что задачи с указанным ID нет.
var ErrTaskNotFound = errors.New("task not found")

// ErrConflict сообщает, что задача уже изменена другим запросом.
var ErrConflict = errors.New("task has changed; refresh and try again")

// ErrHistoryEmpty сообщает, что для отмены нет сохранённых действий.
var ErrHistoryEmpty = errors.New("history is empty")
