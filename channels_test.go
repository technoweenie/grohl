package grohl

import (
	"testing"
)

func TestLog(t *testing.T) {
	channel := make(chan LogData, 1)
	logger, channel := NewChannelLogger(channel)
	data := LogData{"a": 1}
	logger.Log(data)

	recv := <-channel

	if recvKeys := len(recv); recvKeys != len(data) {
		t.Errorf("Wrong number of keys: %d (%s)", recvKeys, recv)
	}

	if data["a"] != recv["a"] {
		t.Errorf("Received: %s", recv)
	}
}
