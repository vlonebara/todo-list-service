package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// loadEnvFile загружает переменные из файла, не заменяя значения окружения.
func loadEnvFile(file string) error {
	input, err := os.Open(file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if errutil.CheckError(err) {
		return err
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if !ok || key == "" {
			return fmt.Errorf("invalid .env entry on line %d", lineNumber)
		}
		if strings.HasPrefix(value, `"`) {
			value, err = strconv.Unquote(value)
			if errutil.CheckError(err) {
				return fmt.Errorf("invalid value on .env line %d", lineNumber)
			}
		} else if len(value) >= 2 && strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			value = value[1 : len(value)-1]
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); errutil.CheckError(err) {
			return fmt.Errorf("set environment variable %s: %w", key, err)
		}
	}
	return scanner.Err()
}
