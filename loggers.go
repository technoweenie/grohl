package grohl

import (
	"fmt"
	"io"
)

// A really basic logger that builds lines and writes to any io.Writer.  This
// expects the writers to be threadsafe.
type IoLogger struct {
	stream  io.Writer
	AddTime bool
}

func (l *IoLogger) Log(data Data) {
	fullLine := fmt.Sprintf("%s\n", BuildLog(data, l.AddTime))
	l.stream.Write([]byte(fullLine))
}

type ChannelLogger struct {
	channel chan Data
}

func (l *ChannelLogger) Log(data Data) {
	l.channel <- data
}
