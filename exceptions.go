package grohl

import (
	"bytes"
	"reflect"
	"runtime/debug"
	"strconv"
)

type ExceptionReporter interface {
	Report(err error, data Data) error
}

// Implementation of ExceptionReporter that writes to a grohl logger.
func (c *Context) Report(err error, data Data) error {
	merged := c.Merge(data)
	errorToMap(err, merged)

	if c.ExceptionReporter != nil {
		return c.ExceptionReporter.Report(err, merged)
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

func ErrorBacktrace(err error) string {
	lines := errorBacktraceBytes(err)
	return string(bytes.Join(lines, byteLineBreak))
}

func ErrorBacktraceLines(err error) []string {
	byteLines := errorBacktraceBytes(err)
	lines := make([]string, len(byteLines))
	for i, byteline := range byteLines {
		lines[i] = string(byteline)
	}
	return lines
}

func errorBacktraceBytes(err error) [][]byte {
	backtrace := debug.Stack()
	all := bytes.Split(backtrace, byteLineBreak)
	return all[4 : len(all)-1]
}

func ErrorId(err error) string {
	id := int(reflect.ValueOf(err).Pointer())
	return strconv.Itoa(id)
}

func errorToMap(err error, data Data) {
	data["at"] = "exception"
	data["class"] = reflect.TypeOf(err).String()
	data["message"] = err.Error()
	data["exception_id"] = ErrorId(err)
}

var byteLineBreak = []byte{'\n'}
