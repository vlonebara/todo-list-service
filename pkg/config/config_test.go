package config

import (
	"testing"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// TestLoad проверяет локальный запуск и обязательный пароль для внешней сети.
func TestLoad(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, item := range []struct {
		name, host, password, port string
		valid, secure              bool
	}{
		{"local without password", "127.0.0.1", "", "7540", true, false},
		{"public without password", "0.0.0.0", "", "7540", false, false},
		{"public with password", "0.0.0.0", "test-password", "7540", true, true},
		{"invalid port", "127.0.0.1", "", "0", false, false},
	} {
		t.Run(item.name, func(t *testing.T) {
			t.Setenv("TODO_HOST", item.host)
			t.Setenv("TODO_PASSWORD", item.password)
			t.Setenv("TODO_PORT", item.port)
			t.Setenv("TODO_COOKIE_SECURE", "")
			cfg, err := Load()
			if (!errutil.CheckError(err)) != item.valid {
				t.Fatalf("unexpected result: %v", err)
			}
			if item.valid && cfg.SecureCookies != item.secure {
				t.Fatal("unexpected cookie mode")
			}
		})
	}
}
