package restartableserver

import (
	"fmt"
	"os"
	"time"
)

// Formats Duration as hh:mm:ss
func FormatDuration(d time.Duration) string {
	r := d.Round(time.Second)

	h := r / time.Hour
	r -= h * time.Hour
	m := r / time.Minute
	r -= m * time.Minute
	s := r / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// Retrieve environment variable with fallback to default value if not present
func Getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
