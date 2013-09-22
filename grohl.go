package grohl

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type Data map[string]interface{}

type Logger interface {
	Log(Data) error
}

var CurrentLogger Logger = NewIoLogger(nil)
var CurrentContext = &Context{make(Data), CurrentLogger, "s", nil}
var CurrentStatter Statter = CurrentContext

func Log(data Data) {
	CurrentContext.Log(data)
}

func Report(err error, data Data) {
	CurrentContext.Report(err, data)
}

func Counter(sampleRate float32, bucket string, n ...int) {
	CurrentStatter.Counter(sampleRate, bucket, n...)
}

func Timing(sampleRate float32, bucket string, d ...time.Duration) {
	CurrentStatter.Timing(sampleRate, bucket, d...)
}

func Gauge(sampleRate float32, bucket string, value ...string) {
	CurrentStatter.Gauge(sampleRate, bucket, value...)
}

func SetLogger(logger Logger) Logger {
	if logger == nil {
		logger = NewIoLogger(nil)
	}

	CurrentLogger = logger
	CurrentContext.Logger = logger

	return logger
}

func NewChannelLogger(channel chan Data) (*ChannelLogger, chan Data) {
	if channel == nil {
		channel = make(chan Data)
	}
	return &ChannelLogger{channel}, channel
}

func NewIoLogger(stream io.Writer) *IoLogger {
	if stream == nil {
		stream = os.Stdout
	}

	return &IoLogger{stream, true}
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

func WriteBufferedLogs(ch chan Data, filename string, upperbound int) {
	logger := NewBufferedLogger(filename, upperbound)
	for {
		logger.Log(<-ch)
	}
}

func NewContext(data Data) *Context {
	return CurrentContext.New(data)
}

func AddContext(key string, value interface{}) {
	CurrentContext.Add(key, value)
}

func DeleteContext(key string) {
	CurrentContext.Delete(key)
}

func NewTimer(data Data) *Timer {
	return CurrentContext.Timer(data)
}

func SetTimeUnit(unit string) {
	CurrentContext.TimeUnit = unit
}

func TimeUnit() string {
	return CurrentContext.TimeUnit
}
