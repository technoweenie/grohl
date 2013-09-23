package grohl

import (
	"bytes"
	"strings"
	"testing"
)

func TestIoLog(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := NewIoLogger(buf)
	logger.AddTime = false
	logger.Log(Data{"a": 1})
	expected := "a=1\n"

	if actual := buf.String(); actual != expected {
		t.Errorf("e: %s\na: %s", expected, actual)
	}
}

func TestChannelLog(t *testing.T) {
	channel := make(chan Data, 1)
	logger, channel := NewChannelLogger(channel)
	data := Data{"a": 1}
	logger.Log(data)

	recv := <-channel

	if recvKeys := len(recv); recvKeys != len(data) {
		t.Errorf("Wrong number of keys: %d (%s)", recvKeys, recv)
	}

	if data["a"] != recv["a"] {
		t.Errorf("Received: %s", recv)
	}
}

func setupLogger() (*Context, chan Data) {
	ch := make(chan Data, 100)
	logger, _ := NewChannelLogger(ch)
	context := NewContext(nil)
	context.Logger = logger
	return context, ch
}

func logged(ch chan Data) string {
	close(ch)
	lines := make([]string, len(ch))
	i := 0

	for data := range ch {
		lines[i] = BuildLog(data, false)
		i = i + 1
	}

	return strings.Join(lines, "\n")
}
