package service

import "testing"

// TestParseDate проверяет единые правила дат для всех входов приложения.
func TestParseDate(t *testing.T) {
	for _, item := range []struct {
		name, value string
		valid       bool
	}{
		{"valid date", "20240229", true},
		{"last supported year", "99991231", true},
		{"missing date", "", false},
		{"short date", "2024011", false},
		{"zero year", "00000101", false},
		{"invalid day", "20240230", false},
		{"wrong format", "29.02.2024", false},
	} {
		t.Run(item.name, func(t *testing.T) {
			_, err := ParseDate(item.value)
			if (err == nil) != item.valid {
				t.Fatalf("unexpected result: %v", err)
			}
		})
	}
}
