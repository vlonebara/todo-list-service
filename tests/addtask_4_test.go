package tests

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type task struct {
	date    string
	title   string
	comment string
	repeat  string
}

// fields формирует одинаковые JSON-поля для создания и редактирования задачи.
func (task task) fields() map[string]any {
	return map[string]any{
		"date": task.date, "title": task.title,
		"comment": task.comment, "repeat": task.repeat,
	}
}

func TestAddTask(t *testing.T) {
	db := openDB(t)

	tbl := []task{
		{"20240129", "", "", ""},
		{"20240192", "Qwerty", "", ""},
		{"28.01.2024", "Заголовок", "", ""},
		{"20240112", "Заголовок", "", "w"},
		{"20240212", "Заголовок", "", "ooops"},
	}
	for _, v := range tbl {
		m, err := postJSON("api/task", v.fields(), http.MethodPost)
		assert.NoError(t, err)

		e, ok := m["error"]
		assert.False(t, !ok || len(fmt.Sprint(e)) == 0,
			"Ожидается ошибка для задачи %v", v)
	}

	now := time.Now()

	check := func() {
		for _, v := range tbl {
			today := v.date == "today"
			if today {
				v.date = now.Format(`20060102`)
			}
			m, err := postJSON("api/task", v.fields(), http.MethodPost)
			assert.NoError(t, err)

			e, ok := m["error"]
			if ok && len(fmt.Sprint(e)) > 0 {
				t.Errorf("Неожиданная ошибка %v для задачи %v", e, v)
				continue
			}
			var mid any
			mid, ok = m["id"]
			if !ok {
				t.Errorf("Не возвращён id для задачи %v", v)
				continue
			}
			id := fmt.Sprint(mid)

			task := getDBTask(t, db, id)
			assert.Equal(t, id, strconv.FormatInt(task.ID, 10))

			assert.Equal(t, v.title, task.Title)
			assert.Equal(t, v.comment, task.Comment)
			assert.Equal(t, v.repeat, task.Repeat)
			if task.Date < now.Format(`20060102`) {
				t.Errorf("Дата не может быть меньше сегодняшней %v", v)
				continue
			}
			if today && task.Date != now.Format(`20060102`) {
				t.Errorf("Дата должна быть сегодняшняя %v", v)
			}
		}
	}

	tbl = []task{
		{"", "Заголовок", "", ""},
		{"20231220", "Сделать что-нибудь", "Хорошо отдохнуть", ""},
		{"20240108", "Уроки", "", "d 10"},
		{"20240102", "Отдых в Сочи", "На лыжах", "y"},
		{"today", "Фитнес", "", "d 1"},
		{"today", "Шмитнес", "", ""},
	}
	check()
	if FullNextDate {
		tbl = []task{
			{"20240129", "Сходить в магазин", "", "w 1,3,5"},
		}
		check()
	}
}
