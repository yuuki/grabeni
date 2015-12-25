package log

import (
	"fmt"
	"os"
)

// Log outputs `format` with `prefix`
func Logf(prefix, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, prefix+": "+format+"\n", args...)
}

func debugf(format string, args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		Logf("debug", format, args...)
	}
}

