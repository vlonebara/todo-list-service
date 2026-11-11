package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/vlonebara/todo-list/pkg/db"
	"github.com/vlonebara/todo-list/pkg/errutil"
	"github.com/vlonebara/todo-list/pkg/handlers"
	"github.com/vlonebara/todo-list/pkg/service"
)

// TestMain запускает все интеграционные тесты с собственной БД и сервером.
func TestMain(m *testing.M) {
	code, err := runTests(m)
	if errutil.CheckError(err) {
		fmt.Fprintln(os.Stderr, "test setup:", err)
		code = 1
	}
	os.Exit(code)
}

// runTests изолирует тестовые данные и закрывает ресурсы после проверок.
func runTests(m *testing.M) (int, error) {
	directory, err := os.MkdirTemp("", "todo-integration-")
	if errutil.CheckError(err) {
		return 1, err
	}
	defer os.RemoveAll(directory)
	DBFile = filepath.Join(directory, "test.db")
	database, err := db.OpenDatabase(DBFile)
	if errutil.CheckError(err) {
		return 1, err
	}
	defer database.Close()
	authentication, err := service.NewAuthService("integration-test-password", "")
	if errutil.CheckError(err) {
		return 1, err
	}
	tasks := service.NewTaskService(db.NewTaskStore(database))
	server := httptest.NewServer(handlers.New(tasks, authentication, false, "../web"))
	defer server.Close()
	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if errutil.CheckError(err) {
		return 1, err
	}
	Port, err = strconv.Atoi(port)
	if errutil.CheckError(err) {
		return 1, err
	}
	Token, err = testToken()
	if errutil.CheckError(err) {
		return 1, err
	}
	return m.Run(), nil
}

// testToken получает свежий JWT тем же запросом, которым пользуется фронтенд.
func testToken() (string, error) {
	response, err := sendRequest(http.MethodPost, "api/signin",
		bytes.NewBufferString(`{"password":"integration-test-password"}`), "")
	if errutil.CheckError(err) {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("test signin returned %d", response.StatusCode)
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); errutil.CheckError(err) {
		return "", err
	}
	if result.Token == "" {
		return "", fmt.Errorf("test signin returned an empty token")
	}
	return result.Token, nil
}
