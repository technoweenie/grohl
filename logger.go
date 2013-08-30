package scrolls

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Logger struct {
	Stream io.Writer
}

func NewLogger(stream io.Writer) *Logger {
	if stream == nil {
		stream = os.Stdout
	}
	return &Logger{stream}
}

func (l *Logger) Log(context map[string]interface{}) {
	l.Stream.Write([]byte(l.buildLine(context)))
}

func (l *Logger) buildLine(context map[string]interface{}) string {
	if context != nil {
		return l.convertContext(context)
	} else {
		return empty
	}
}

func (l *Logger) convertContext(context map[string]interface{}) string {
	pieces := make([]string, len(context))
	index := 0

	for key, value := range context {
		pieces[index] = fmt.Sprintf("%s=%s", key, value)
		index = index + 1
	}

	return strings.Join(pieces, space)
}

const (
	space = " "
	empty = ""
)
