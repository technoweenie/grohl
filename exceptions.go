package grohl

import (
	"bytes"
	"reflect"
	"runtime/debug"
	"strconv"
)

type ExceptionReporter interface {
	Report(err error, data map[string]interface{})
}

// Implementation of ExceptionReporter that writes to a grohl logger.
type ExceptionLogReporter struct {
	logger Logger
}

func (r *ExceptionLogReporter) Report(err error, data map[string]interface{}) {
	r.logger.Log(data)
	for _, line := range ErrorBacktraceLines(err) {
		data["site"] = line
		r.logger.Log(data)
	}
}

type ExceptionLogger struct {
	Reporter ExceptionReporter
	*IoLogger
}

func newExceptionLogger(logger *IoLogger, reporter ExceptionReporter) *ExceptionLogger {
	if reporter == nil {
		reporter = &ExceptionLogReporter{logger}
	}

	return &ExceptionLogger{reporter, logger}
}

func (l *ExceptionLogger) Report(err error, data map[string]interface{}) {
	l.ReportException(l.Reporter, err, data)
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

func errorToMap(err error, data map[string]interface{}) {
	data["at"] = "exception"
	data["class"] = reflect.TypeOf(err).String()
	data["message"] = err.Error()
	data["exception_id"] = ErrorId(err)
}

var byteLineBreak = []byte{'\n'}
