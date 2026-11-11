package handlers

import "testing"

// checkError останавливает тест при ошибке и указывает место вызова помощника.
func checkError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
