package scrolls

import (
	"bytes"
	"testing"
)

func TestLogData(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.Log(map[string]interface{}{
		"a": "1", "b": "2",
	})

	if result := buf.String(); result != "a=1 b=2" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestLogDataWithContext(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddContext("a", "1")
	logger.AddContext("b", "1")

	logger.Log(map[string]interface{}{
		"b": "2", "c": "3",
	})

	if result := buf.String(); result != "a=1 b=2 c=3" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestLogEmptyData(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.Log(nil)

	if result := buf.String(); result != "" {
		t.Errorf("Bad log output: %s", result)
	}
}

func loggerWithBuffer() (*Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return NewLogger(buf), buf
}
