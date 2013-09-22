package grohl

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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

type BufferedLogger struct {
	Filename   string
	Position   int
	Upperbound int
	AddTime    bool
	Fallback   io.WriteCloser
	buffer     []byte
	isFile     bool
}

func NewBufferedLogger(filename string, upperbound int) *BufferedLogger {
	isFile := len(filename) > 0

	if isFile {
		os.MkdirAll(filepath.Dir(filename), 0744)
	}

	return &BufferedLogger{
		Filename:   filename,
		Position:   0,
		Upperbound: upperbound,
		AddTime:    true,
		Fallback:   &nopcloser{os.Stdout},
		buffer:     make([]byte, upperbound),
		isFile:     isFile,
	}
}

func WriteBufferedLogs(logch chan Data, filename string, upperbound int) {
	logger := NewBufferedLogger(filename, upperbound)
	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-sigch:
			logger.Flush()
		case data := <-logch:
			logger.Log(data)
		}
	}
}

func (l *BufferedLogger) Log(data Data) {
	line := l.BuildLog(data)
	length := len(line)

	if length > l.Upperbound {
		l.Write(line)
		return
	}

	if (length + l.Position) > l.Upperbound {
		l.Flush()
	}

	copy(l.buffer[l.Position:], line)
	l.Position += length
}

func (l *BufferedLogger) Write(data []byte) (int, error) {
	stream := l.openStream()
	written, err := stream.Write(l.buffer[:l.Position])
	stream.Close()
	return written, err
}

func (l *BufferedLogger) BuildLog(data Data) []byte {
	return []byte(fmt.Sprintf("%s\n", BuildLog(data, l.AddTime)))
}

func (l *BufferedLogger) Flush() {
	if l.Position == 0 {
		return
	}

	l.Write(l.buffer[:l.Position])
	l.Position = 0
}

func (l *BufferedLogger) openStream() io.WriteCloser {
	if l.isFile {
		file, err := os.OpenFile(l.Filename, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
		if err == nil {
			return file
		} else {
			fmt.Printf("Error opening log file %s: %s\n", l.Filename, err.Error())
		}
	}

	return l.Fallback
}

type nopcloser struct {
	io.Writer
}

func (c *nopcloser) Close() error {
	return nil
}
