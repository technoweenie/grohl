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

func (l *IoLogger) Log(data Data) error {
	fullLine := fmt.Sprintf("%s\n", BuildLog(data, l.AddTime))
	_, err := l.stream.Write([]byte(fullLine))
	return err
}

type ChannelLogger struct {
	channel chan Data
}

func (l *ChannelLogger) Log(data Data) error {
	l.channel <- data
	return nil
}
