package scrolls

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Logger struct {
	Stream  io.Writer
	context map[string]interface{}
}

func NewLogger(stream io.Writer) *Logger {
	return NewLoggerWithContext(stream, make(map[string]interface{}))
}

func NewLoggerWithContext(stream io.Writer, context map[string]interface{}) *Logger {
	if stream == nil {
		stream = os.Stdout
	}
	return &Logger{stream, context}
}

func (l *Logger) Log(data map[string]interface{}) {
	l.Stream.Write([]byte(l.buildLine(data)))
}

func (l *Logger) NewContext(data map[string]interface{}) *Logger {
	return NewLoggerWithContext(l.Stream, dupeMaps(l.context, data))
}

func (l *Logger) AddContext(key string, value interface{}) {
	l.context[key] = value
}

func (l *Logger) DeleteContext(key string) {
	delete(l.context, key)
}

func (l *Logger) buildLine(data map[string]interface{}) string {
	merged := dupeMaps(l.context, data)
	pieces := make([]string, len(merged))
	l.convertDataMap(merged, pieces)
	return strings.Join(pieces, space)
}

func (l *Logger) convertDataMap(data map[string]interface{}, pieces []string) {
	index := 0
	for key, value := range data {
		pieces[index] = fmt.Sprintf("%s=%s", key, value)
		index = index + 1
	}
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
	space = " "
	empty = ""
)
