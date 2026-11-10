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
	"time"

	"github.com/vlonebara/todo-list/pkg/db"
	"github.com/vlonebara/todo-list/pkg/handlers"
	"github.com/vlonebara/todo-list/pkg/service"
)

// TestMain запускает все интеграционные тесты с собственной БД и сервером.
func TestMain(m *testing.M) {
	code, err := runTests(m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "test setup:", err)
		code = 1
	}
	os.Exit(code)
}

// runTests изолирует тестовые данные и закрывает ресурсы после проверок.
func runTests(m *testing.M) (int, error) {
	directory, err := os.MkdirTemp("", "todo-integration-")
	if err != nil {
		return 1, err
	}
	defer os.RemoveAll(directory)
	DBFile = filepath.Join(directory, "test.db")
	database, err := db.OpenDatabase(DBFile)
	if err != nil {
		return 1, err
	}
	defer database.Close()
	authentication, err := service.NewAuthService("integration-test-password", "")
	if err != nil {
		return 1, err
	}
	tasks := service.NewTaskService(db.NewTaskStore(database))
	server := httptest.NewServer(handlers.New(tasks, authentication, false, "../web"))
	defer server.Close()
	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		return 1, err
	}
	Port, err = strconv.Atoi(port)
	if err != nil {
		return 1, err
	}
	if err := os.Setenv("TODO_PORT", port); err != nil {
		return 1, err
	}
	if err := os.Setenv("TODO_DBFILE", DBFile); err != nil {
		return 1, err
	}
	Token, err = testToken(server.URL)
	if err != nil {
		return 1, err
	}
	return m.Run(), nil
}

// testToken получает свежий JWT тем же запросом, которым пользуется фронтенд.
func testToken(address string) (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Post(address+"/api/signin", "application/json",
		bytes.NewBufferString(`{"password":"integration-test-password"}`))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("test signin returned %d", response.StatusCode)
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Token == "" {
		return "", fmt.Errorf("test signin returned an empty token")
	}
	return result.Token, nil
}
