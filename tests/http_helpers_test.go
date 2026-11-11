package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// sendRequest отправляет запрос с общими заголовками, таймаутом и выбранным токеном.
func sendRequest(method, path string, body io.Reader, token string) (*http.Response, error) {
	request, err := http.NewRequest(method, getURL(path), body)
	if errutil.CheckError(err) {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.AddCookie(&http.Cookie{Name: "token", Value: token})
	}
	client := http.Client{Timeout: 5 * time.Second}
	return client.Do(request)
}

// requestJSON отправляет поля задачи и читает ответ, закрывая его тело.
func requestJSON(path string, values map[string]any, method string) ([]byte, error) {
	var data []byte
	if len(values) > 0 {
		var err error
		data, err = json.Marshal(values)
		if errutil.CheckError(err) {
			return nil, err
		}
	}
	response, err := sendRequest(method, path, bytes.NewReader(data), Token)
	if errutil.CheckError(err) {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

// readJSON получает ответ API и разбирает его в указанную структуру.
func readJSON(path string, values map[string]any, method string, result any) error {
	body, err := requestJSON(path, values, method)
	if errutil.CheckError(err) {
		return err
	}
	return json.Unmarshal(body, result)
}

// postJSON возвращает JSON-объект API для проверки полей id и error.
func postJSON(path string, values map[string]any, method string) (map[string]any, error) {
	var result map[string]any
	err := readJSON(path, values, method, &result)
	return result, err
}
