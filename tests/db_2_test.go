package tests

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

// count возвращает количество задач для проверки изменений таблицы.
func count(db *sqlx.DB) (int, error) {
	var count int
	err := db.Get(&count, `SELECT count(id) FROM scheduler`)
	return count, err
}

// openDB открывает тестовую БД и закрывает её после сценария.
func openDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Connect("sqlite", DBFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}

// getDBTask читает сохранённую задачу для сравнения с отправленными данными.
func getDBTask(t *testing.T, db *sqlx.DB, id any) Task {
	t.Helper()
	var task Task
	if err := db.Get(&task, `SELECT * FROM scheduler WHERE id=?`, id); err != nil {
		t.Fatal(err)
	}
	return task
}

func TestDB(t *testing.T) {
	db := openDB(t)

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
	VALUES (?, 'Todo', 'Комментарий', '')`, today)
	assert.NoError(t, err)

	id, err := res.LastInsertId()
	assert.NoError(t, err)

	task := getDBTask(t, db, id)
	assert.Equal(t, id, task.ID)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Комментарий`, task.Comment)

	_, err = db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
