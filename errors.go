package grohl

import (
	"bytes"
	"reflect"
	"runtime/debug"
)

type ErrorReporter interface {
	Report(err error, data Data) error
}

// Report writes the error to the ErrorReporter, or logs it if there is none.
func (c *Context) Report(err error, data Data) error {
	merged := c.Merge(data)
	errorToMap(err, merged)

	if c.ErrorReporter != nil {
		return c.ErrorReporter.Report(err, merged)
	} else {
		var logErr error
		logErr = c.log(merged)
		if logErr != nil {
			return logErr
		}

		for _, line := range ErrorBacktraceLines(err) {
			lineData := dupeMaps(merged)
			lineData["site"] = line
			logErr = c.log(lineData)
			if logErr != nil {
				return logErr
			}
		}
		return nil
	}
}

// ErrorBacktrace creates a backtrace of the call stack.
func ErrorBacktrace(err error) string {
	return string(debug.Stack())
}

// ErrorBacktraceLines creates a backtrace of the call stack, split into lines.
func ErrorBacktraceLines(err error) []string {
	byteLines := errorBacktraceBytes(err)
	lines := make([]string, 0, len(byteLines))

	// skip top two frames which are this method and `errorBacktraceBytes`
	for i := 2; i < len(byteLines); i++ {
		lines = append(lines, string(byteLines[i]))
	}
	return lines
}

func errorBacktraceBytes(err error) [][]byte {
	return bytes.Split(debug.Stack(), byteLineBreak)
}

func errorToMap(err error, data Data) {
	data["at"] = "exception"
	data["class"] = reflect.TypeOf(err).String()
	data["message"] = err.Error()
}

var byteLineBreak = []byte{'\n'}
