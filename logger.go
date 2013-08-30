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
	if stream == nil {
		stream = os.Stdout
	}
	return &Logger{stream, make(map[string]interface{})}
}

func (l *Logger) Log(data map[string]interface{}) {
	l.Stream.Write([]byte(l.buildLine(data)))
}

func (l *Logger) AddContext(key string, value interface{}) {
	l.context[key] = value
}

func (l *Logger) DeleteContext(key string) {
	delete(l.context, key)
}

func (l *Logger) buildLine(data map[string]interface{}) string {
	merged := make(map[string]interface{})
	mergeMaps(l.context, merged)
	mergeMaps(data, merged)

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

func mergeMaps(original map[string]interface{}, copy map[string]interface{}) {
	if original == nil {
		return
	}

	for key, value := range original {
		copy[key] = value
	}
}

const (
	space = " "
	empty = ""
)
