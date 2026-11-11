package db

import (
	"fmt"
	"testing"

	"github.com/vlonebara/todo-list/pkg/model"
)

// TestTaskList проверяет страницы, русские буквы и буквальный поиск символов.
func TestTaskList(t *testing.T) {
	store := testStore(t)
	for i := 0; i < 20; i++ {
		addTestTask(t, store, model.Task{Date: "20300101", Title: fmt.Sprintf("Бассейн %02d", i), Comment: "C++ A&B #tag 100%"})
	}
	cases := []struct {
		name, search string
		offset, want int
	}{
		{"first page", "", 0, 15},
		{"second page", "", 15, 5},
		{"lowercase Russian", "бассейн", 0, 15},
		{"search second page", "БАССЕЙН", 15, 5},
		{"plus signs", "c++", 0, 15},
		{"ampersand", "a&b", 0, 15},
		{"fragment character", "#tag", 0, 15},
		{"literal percent", "100%", 0, 15},
		{"no wildcard search", "missing%", 0, 0},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			tasks, err := store.List(t.Context(), model.TaskFilter{Text: item.search, Offset: item.offset})
			if err != nil || len(tasks) != item.want {
				t.Fatalf("got %d tasks, %v; want %d", len(tasks), err, item.want)
			}
		})
	}
	all, err := store.All(t.Context())
	if err != nil || len(all) != 20 {
		t.Fatalf("incomplete export: %d tasks, %v", len(all), err)
	}
}

// TestTaskDateFilters проверяет календарные границы выборки.
func TestTaskDateFilters(t *testing.T) {
	store := testStore(t)
	for _, date := range []string{"20290101", "20300101", "20310101"} {
		addTestTask(t, store, model.Task{Date: date, Title: date})
	}
	for _, item := range []struct {
		name   string
		filter model.TaskFilter
		want   int
	}{
		{"exact date", model.TaskFilter{Date: "20300101"}, 1},
		{"older tasks", model.TaskFilter{BeforeDate: "20300101"}, 1},
		{"future tasks", model.TaskFilter{FromDate: "20300102"}, 1},
		{"one day", model.TaskFilter{FromDate: "20300101", BeforeDate: "20300102"}, 1},
	} {
		t.Run(item.name, func(t *testing.T) {
			tasks, err := store.List(t.Context(), item.filter)
			if err != nil || len(tasks) != item.want {
				t.Fatalf("unexpected selection: %v, %v", tasks, err)
			}
		})
	}
}
