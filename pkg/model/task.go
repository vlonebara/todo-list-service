package model

// Task содержит данные задачи и не зависит от БД или HTTP.
type Task struct {
	ID      int64
	Date    string
	Title   string
	Comment string
	Repeat  string
}

// TaskFilter задаёт поиск задач по тексту или дате.
type TaskFilter struct {
	Text string
	Date string
}

// TaskChange хранит состояние задачи до и после действия.
type TaskChange struct {
	ID        int64
	Action    string
	Before    *Task
	After     *Task
	CreatedAt string
	Undone    bool
}
