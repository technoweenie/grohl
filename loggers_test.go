package grohl

import (
	"bytes"
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

type loggerBuffer struct {
	channel chan Data
	t       *testing.T
	lines   []builtLogLine
	index   int
}

func (b *loggerBuffer) Lines() []builtLogLine {
	if b.lines == nil {
		close(b.channel)
		b.lines = make([]builtLogLine, len(b.channel))
		i := 0

		for data := range b.channel {
			b.lines[i] = buildLogLine(data)
			i = i + 1
		}
	}

	return b.lines
}

func (b *loggerBuffer) AssertLine(parts ...string) {
	lines := b.Lines()
	if b.index < 0 || b.index >= len(lines) {
		b.t.Errorf("No line %d", b.index)
		return
	}

	AssertBuiltLine(b.t, lines[b.index], parts...)
	b.index += 1
}

func (b *loggerBuffer) AssertEOF() {
	lines := b.Lines()
	if b.index < 0 {
		b.t.Errorf("Invalid index %d", b.index)
		return
	}

	if b.index < len(lines) {
		b.t.Errorf("Not EOF, on line %d", b.index)
		return
	}
}

func setupLogger(t *testing.T) (*Context, *loggerBuffer) {
	ch := make(chan Data, 100)
	logger, _ := NewChannelLogger(ch)
	context := NewContext(nil)
	context.Logger = logger
	return context, &loggerBuffer{channel: ch, t: t}
}
