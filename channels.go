package grohl

type ChannelLogger struct {
	channel chan LogData
}

func SetupChannelLogger(channel chan LogData) chan LogData {
  logger, channel := NewChannelLogger(channel)
  DefaultLogger = logger
  return channel
}

func NewChannelLogger(channel chan LogData) (*ChannelLogger, chan LogData) {
	if channel == nil {
		channel = make(chan LogData)
	}
	return &ChannelLogger{channel}, channel
}

func (l *ChannelLogger) Log(data LogData) {
	l.channel <- data
}
