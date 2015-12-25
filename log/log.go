package log

import (
	"fmt"
	"os"
)

var IsDebug = false

// Log outputs `format` with `prefix`
func Logf(prefix, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, prefix+": "+format+"\n", args...)
}

func Debugf(format string, args ...interface{}) {
	if IsDebug == true {
		Logf("debug", format, args...)
	}
}

