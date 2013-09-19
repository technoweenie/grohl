package grohl

import (
	"io"
	"os"
	"time"
)

type Data map[string]interface{}

type Logger interface {
	Log(Data)
}

var CurrentLogger Logger = NewIoLogger(nil)
var CurrentContext = &Context{make(Data), CurrentLogger, "s"}
var CurrentReporter ExceptionReporter = CurrentContext
var CurrentStatter Statter = CurrentContext

func Log(data Data) {
	CurrentContext.Log(data)
}

func Report(err error, data Data) {
	CurrentReporter.Report(err, data)
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
