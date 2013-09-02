package grohl

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type IoLogger struct {
	TimeUnit string
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

	return &IoLogger{defaultTimeUnit, stream, context}
}

func (l *IoLogger) Log(data map[string]interface{}) {
	fullLine := fmt.Sprintf("%s\n", l.buildLine(data))
	l.stream.Write([]byte(fullLine))
}

func (l *IoLogger) NewContext(data map[string]interface{}) *IoLogger {
	ctx := NewLoggerWithContext(l.stream, dupeMaps(l.context, data))
	ctx.TimeUnit = l.TimeUnit
	return ctx
}

func (l *IoLogger) AddContext(key string, value interface{}) {
	l.context[key] = value
}

func (l *IoLogger) DeleteContext(key string) {
	delete(l.context, key)
}

func (l *IoLogger) buildLine(data map[string]interface{}) string {
	merged := dupeMaps(l.context, data)
	pieces := make([]string, len(merged))

	index := 0
	for key, value := range merged {
		pieces[index] = fmt.Sprintf("%s=%s", key, value)
		index = index + 1
	}

	return strings.Join(pieces, space)
}

func dupeMaps(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, orig := range maps {
		for key, value := range orig {
			merged[key] = value
		}
	}
	return merged
}

const (
	space           = " "
	defaultTimeUnit = "s"
)
