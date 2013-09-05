package grohl

import (
	"fmt"
	"io"
	"os"
)

// A really basic logger that builds lines and writes to any io.Writer.  This
// expects the writers to be threadsafe.  If they're not, wrap them with
// SyncWriter() to protect Write() calls with a mutex.
type IoLogger struct {
	TimeUnit string
	AddTime  bool
	stream   io.Writer
	context  map[string]interface{}
}

type LogData map[string]interface{}

type Logger interface {
	Log(data map[string]interface{})
}

func NewLogger(stream io.Writer) *IoLogger {
	return NewLoggerWithContext(stream, nil)
}

func NewLoggerWithContext(stream io.Writer, context map[string]interface{}) *IoLogger {
	if stream == nil {
		stream = SyncWriter(os.Stdout)
	}

	if context == nil {
		context = make(map[string]interface{})
	}

	return &IoLogger{defaultTimeUnit, true, stream, context}
}

func (l *IoLogger) Log(data map[string]interface{}) {
	fullLine := fmt.Sprintf("%s\n", buildLine(l.context, data, l.AddTime))
	l.stream.Write([]byte(fullLine))
}

func (l *IoLogger) NewContext(data map[string]interface{}) *IoLogger {
	ctx := NewLoggerWithContext(l.stream, dupeMaps(l.context, data))
	ctx.TimeUnit = l.TimeUnit
	ctx.AddTime = l.AddTime
	return ctx
}

func (l *IoLogger) NewExceptionLogger(reporter ExceptionReporter) *ExceptionLogger {
	return newExceptionLogger(l, reporter)
}

func (l *IoLogger) AddContext(key string, value interface{}) {
	l.context[key] = value
}

func (l *IoLogger) DeleteContext(key string) {
	delete(l.context, key)
}

const defaultTimeUnit = "s"
