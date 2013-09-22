package grohl

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
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
	Position   int
	Upperbound int
	AddTime    bool
	Timeout    time.Duration
	Writer     io.WriteCloser
	buffer     []byte
}

func NewBufferedLogger(filename string, upperbound int) *BufferedLogger {
	stdout := &nopcloser{os.Stdout}
	var writer io.WriteCloser = stdout

	if len(filename) > 0 {
		err := os.MkdirAll(filepath.Dir(filename), 0744)
		if err == nil {
			writer = &bufferedFileWriter{filename, stdout}
		} else {
			fmt.Printf("Error opening log file %s: %s\n", filename, err.Error())
		}
	}

	return &BufferedLogger{
		Position:   0,
		Upperbound: upperbound,
		AddTime:    true,
		Timeout:    1 * time.Minute,
		Writer:     writer,
		buffer:     make([]byte, upperbound),
	}
}

func WriteBufferedLogs(logch chan Data, filename string, upperbound int) {
	NewBufferedLogger(filename, upperbound).Watch(logch)
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
	return l.Writer.Write(data)
}

func (l *BufferedLogger) Watch(logch chan Data) {
	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-sigch:
			l.Flush()
		case data := <-logch:
			if data != nil {
				l.Log(data)
			}
		case <-time.After(l.Timeout):
			l.Flush()
		}
	}
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

type bufferedFileWriter struct {
	Filename string
	Fallback io.WriteCloser
}

func (w *bufferedFileWriter) Write(data []byte) (int, error) {
	stream := w.openStream()
	defer stream.Close()
	return stream.Write(data)
}

func (w *bufferedFileWriter) Close() error {
	return nil
}

func (w *bufferedFileWriter) openStream() io.WriteCloser {
	file, err := os.OpenFile(w.Filename, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
	if err == nil {
		return file
	} else {
		fmt.Printf("Error opening log file %s: %s\n", w.Filename, err.Error())
	}

	return w.Fallback
}

type nopcloser struct {
	io.Writer
}

func (c *nopcloser) Close() error {
	return nil
}
