package grohl

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestLogData(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.Log(LogData{
		"a": "1", "b": "2",
	})

	if result := logged(buf); result != "a=1 b=2" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestLogDataWithTime(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddTime = true
	logger.Log(LogData{
		"a": "1", "b": "2",
	})

	result := logged(buf)
	pieces := strings.Split(result, " ")

	if piecesLen := len(pieces); piecesLen != 3 {
		t.Fatalf("Wrong number of log keys: %d", piecesLen)
	}

	key := pieces[0]
	matched, err := regexp.MatchString(`^now=\d{4}\-\d{2}\-\d{2}T\d{2}:\d{2}:\d{2}\+0000$`, key)
	if err != nil || !matched {
		t.Errorf("Bad first key: %s", key)
	}

	if key = pieces[1]; key != "a=1" {
		t.Errorf("Bad second key: %s", key)
	}

	if key = pieces[2]; key != "b=2" {
		t.Errorf("Bad third key: %s", key)
	}
}

func TestContextObject(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddContext("a", "1")
	logger.AddContext("b", "1")
	context := logger.NewContext(LogData{
		"b": "2", "c": "2",
	})

	context.Log(LogData{
		"c": "3", "d": "4",
	})

	if result := logged(buf); result != "a=1 b=2 c=3 d=4" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestLogDataWithContext(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddContext("a", "1")
	logger.AddContext("b", "1")

	logger.Log(LogData{
		"b": "2", "c": "3",
	})

	if result := logged(buf); result != "a=1 b=2 c=3" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestContextDelete(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddContext("a", "1")
	logger.AddContext("b", "1")
	logger.DeleteContext("b")

	logger.Log(nil)

	if result := logged(buf); result != "a=1" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestLogEmptyData(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.Log(nil)

	if result := logged(buf); result != "" {
		t.Errorf("Bad log output: %s", result)
	}
}

func loggerWithBuffer() (*IoLogger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	logger := NewLogger(buf)
	logger.AddTime = false
	return logger, buf
}

func logged(buf *bytes.Buffer) string {
	return strings.Trim(buf.String(), "\n")
}
