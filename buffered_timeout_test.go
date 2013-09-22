// +build !race

package grohl

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func TestBufferTimeout(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := NewBufferedLogger("", 10)
	logger.AddTime = false
	logger.Writer = &nopcloser{buf}
	logger.Timeout = 1 * time.Millisecond

	ch := make(chan Data, 3)
	ch <- Data{"a": 1}

	go logger.Watch(ch)

	time.Sleep(10 * time.Millisecond)
	if logger.Position != 0 {
		t.Errorf("Not flushed: %d", logger.Position)
	}
	ch <- nil

	expected := "a=1\n"
	if actual := buf.String(); actual != expected {
		t.Errorf("e: %s\na: %s", strconv.Quote(expected), strconv.Quote(actual))
	}
}
