package grohl

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
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

func TestBufferedLog(t *testing.T) {
	file, err := ioutil.TempFile("", "testing")
	if err != nil {
		t.Fatalf("Error opening tempfile: %s", err.Error())
	}

	defer file.Close()
	defer os.Remove(file.Name())

	logger := NewBufferedLogger(file.Name(), 10)
	logger.AddTime = false

	logger.Log(Data{"a": 1})
	logger.Log(Data{"b": 2})

	readFile := func() string {
		by, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatalf("Error reading tempfile: %s", err.Error())
			return ""
		}
		return string(by)
	}

	if actual := readFile(); len(actual) > 0 {
		t.Errorf("expected empty string, got %s", actual)
	}

	logger.Log(Data{"c": 3})

	expected := "a=1\nb=2\n"
	if actual := readFile(); actual != expected {
		t.Errorf("e: %s\na: %s", strconv.Quote(expected), strconv.Quote(actual))
	}

	expected = "c=3\n"
	if actual := string(logger.buffer[:logger.Position]); actual != expected {
		t.Errorf("e: %s\na: %s", strconv.Quote(expected), strconv.Quote(actual))
	}
}

func TestBufferedLogFallback(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := NewBufferedLogger("", 10)
	logger.AddTime = false
	logger.Writer = &nopcloser{buf}

	logger.Log(Data{"a": 1})
	logger.Log(Data{"b": 2})

	if actual := buf.String(); len(actual) > 0 {
		t.Errorf("expected empty string, got %s", actual)
	}

	logger.Log(Data{"c": 3})

	expected := "a=1\nb=2\n"
	if actual := buf.String(); actual != expected {
		t.Errorf("e: %s\na: %s", strconv.Quote(expected), strconv.Quote(actual))
	}

	expected = "c=3\n"
	if actual := string(logger.buffer[:logger.Position]); actual != expected {
		t.Errorf("e: %s\na: %s", strconv.Quote(expected), strconv.Quote(actual))
	}
}

func TestBufferLongLine(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := NewBufferedLogger("", 1)
	logger.AddTime = false
	logger.Writer = &nopcloser{buf}

	logger.Log(Data{"a": 1})
	expected := "a=1\n"
	if actual := buf.String(); actual != expected {
		t.Errorf("e: %s\na: %s", strconv.Quote(expected), strconv.Quote(actual))
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
