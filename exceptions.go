package grohl

import (
	"reflect"
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
}

type ExceptionLogger struct {
	Reporter ExceptionReporter
	*IoLogger
}

func NewExceptionLogger(logger *IoLogger, reporter ExceptionReporter) *ExceptionLogger {
	if logger == nil {
		logger = defaultLogger
	}

	if reporter == nil {
		reporter = &ExceptionLogReporter{logger}
	}

	return &ExceptionLogger{reporter, logger}
}

func (l *ExceptionLogger) Report(err error, data map[string]interface{}) {
	merged := dupeMaps(l.context, data)
	errorToMap(err, merged)
	l.Reporter.Report(err, merged)
}

func errorToMap(err error, data map[string]interface{}) {
	data["at"] = "exception"
	data["class"] = reflect.TypeOf(err).String()
	data["message"] = err.Error()
	data["exception_id"] = int64(reflect.ValueOf(err).Pointer())
}
