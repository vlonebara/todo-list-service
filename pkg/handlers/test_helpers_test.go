package handlers

import (
	"testing"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// checkError останавливает тест при ошибке и указывает место вызова помощника.
func checkError(t *testing.T, err error) {
	t.Helper()
	if errutil.CheckError(err) {
		t.Fatal(err)
	}
}
