package service

import (
	"testing"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// TestNextDate проверяет календарные границы и давние исходные даты.
func TestNextDate(t *testing.T) {
	cases := []struct {
		name, now, start, repeat, want string
	}{
		{"old daily task", "20240126", "00010101", "d 1", "20240127"},
		{"weekly interval", "20240126", "20240113", "d 7", "20240127"},
		{"future task", "20240126", "20240201", "d 7", "20240208"},
		{"leap day", "20240126", "20240229", "y", "20250301"},
		{"old leap day", "20240229", "20000229", "y", "20240301"},
		{"monthly leap day", "20250301", "20240101", "m 29 2", "20280229"},
		{"day overflow", "99991231", "99991231", "d 1", ""},
		{"year overflow", "99991231", "99991231", "y", ""},
		{"weekday overflow", "99991231", "99991231", "w 1", ""},
		{"month overflow", "99991231", "99991231", "m 1", ""},
		{"zero year", "20240126", "00000101", "y", ""},
		{"impossible month", "20240126", "20240101", "m 31 2", ""},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			now, err := time.Parse(dateLayout, item.now)
			if errutil.CheckError(err) {
				t.Fatal(err)
			}
			got, err := NextDate(now, item.start, item.repeat)
			if item.want == "" {
				if !errutil.CheckError(err) {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if errutil.CheckError(err) || got != item.want {
				t.Fatalf("got %q, %v; want %q", got, err, item.want)
			}
		})
	}
}
