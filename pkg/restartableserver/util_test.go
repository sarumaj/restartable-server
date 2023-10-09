package restartableserver

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	for _, tt := range []struct {
		name string
		args time.Duration
		want string
	}{
		{"test#1", 12345689 * time.Millisecond, "03:25:46"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.args)
			if got != tt.want {
				t.Errorf(`FormatDuration(%q) failed: got: %q, want: %q`, tt.args, got, tt.want)
			}
		})
	}
}
