package scrolls

import (
	"bytes"
	"testing"
)

func TestLogToStream(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := NewLogger(buf)
	logger.Log(nil)

	if result := buf.String(); result != "boom" {
		t.Errorf("Bad log output: %s", result)
	}
}
