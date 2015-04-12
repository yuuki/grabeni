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

// ErrorIf outputs log if `err` occurs.
func ErrorIf(err error) bool {
	if err != nil {
		Logf("error", err.Error())
		return true
	}

	return false
}

// DieIf outputs log and exit(1) if `err` occurs.
func DieIf(err error) {
	if err != nil {
		Logf("error", err.Error())
		os.Exit(1)
	}
}

// PanicIf raise panic if `err` occurs.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
