package httputils

import (
	"fmt"
	"os"
)

var debug struct {
	isEnabled     bool
	isAuthEnabled bool
}

// Debug turns on printing every request and response to stderr.
func Debug(enable bool) {
	debug.isEnabled = enable
}

// DebugAuth also prints the Authorization header while Debug is on. Keep it off outside local
// troubleshooting: the header carries the API key and secret.
func DebugAuth(enable bool) {
	debug.isAuthEnabled = enable
}

func debugf(format string, args ...any) {
	if debug.isEnabled {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
