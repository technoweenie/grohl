package grohl

import (
	"fmt"
	"io"
	"os"
)

// A really basic logger that builds lines and writes to any io.Writer.  This
// expects the writers to be threadsafe.
type IoLogger struct {
	stream  io.Writer
	AddTime bool
}

func NewIoLogger(stream io.Writer) *IoLogger {
	if stream == nil {
		stream = os.Stdout
	}

	return &IoLogger{stream, true}
}

func (l *IoLogger) Log(data Data) error {
	line := fmt.Sprintf("%s\n", BuildLog(data, l.AddTime))
	_, err := l.stream.Write([]byte(line))
	return err
}

type ChannelLogger struct {
	channel chan Data
}

func NewChannelLogger(channel chan Data) (*ChannelLogger, chan Data) {
	if channel == nil {
		channel = make(chan Data)
	}
	return &ChannelLogger{channel}, channel
}

func (l *ChannelLogger) Log(data Data) error {
	l.channel <- data
	return nil
}

func Watch(logger Logger, logch chan Data) {
	for {
		data := <-logch
		if data != nil {
			logger.Log(data)
		} else {
			return
		}
	}
}
