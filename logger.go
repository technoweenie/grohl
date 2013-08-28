package scrolls

import (
	"io"
	"os"
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

func (l *Logger) Log(context map[string]string) {
	l.Stream.Write([]byte(l.buildLine(context)))
}

func (l *Logger) buildLine(context map[string]string) string {
	return "boom"
}
