package scrolls

import (
	"bytes"
	"testing"
)

func TestLogContext(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.Log(map[string]interface{}{
		"a": "1", "b": "2",
	})

	if result := buf.String(); result != "a=1 b=2" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestLogEmptyContext(t *testing.T) {
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
