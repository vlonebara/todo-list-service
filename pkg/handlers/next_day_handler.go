package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/service"
)

// nextDayHandler возвращает ближайшую дату задачи по переданному правилу.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if value := r.FormValue("now"); value != "" {
		parsed, err := service.ParseDate(value)
		if errutil.CheckError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	date, err := service.NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if errutil.CheckError(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := w.Write([]byte(date)); errutil.CheckError(err) {
		log.Printf("write next date: %v", err)
	}
}
