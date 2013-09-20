package grohl

import (
	"bytes"
	"reflect"
	"runtime/debug"
	"strconv"
)

type ExceptionReporter interface {
	Report(err error, data Data)
}

// Implementation of ExceptionReporter that writes to a grohl logger.
func (c *Context) Report(err error, data Data) {
	if c.ExceptionReporter != nil {
		c.ExceptionReporter.Report(err, c.Merge(data))
	} else {
		c.report(err, data)
	}
}

func (c *Context) report(err error, data Data) {
	errorToMap(err, data)
	c.Log(data)
	for _, line := range ErrorBacktraceLines(err) {
		data["site"] = line
		c.Log(data)
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
