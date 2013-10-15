package grohl

import (
	"fmt"
	"strings"
	"testing"
)

func TestLogsException(t *testing.T) {
	reporter, buf := setupLogger(t)
	reporter.Add("a", 1)
	reporter.Add("b", 1)

	err := fmt.Errorf("Test")

	reporter.Report(err, Data{"b": 2, "c": 3, "at": "overwrite me"})
	expected := "a=1 b=2 c=3 at=exception class=*errors.errorString message=Test"
	linePrefix := expected + " site="

	for i, line := range strings.Split(buf.String(), "\n") {
		if i == 0 {
			if line != expected {
				t.Errorf("Line does not match:\ne: %s\na: %s", expected, line)
			}
		} else {
			if !strings.HasPrefix(line, linePrefix) {
				t.Errorf("Line %d does not match:\ne: %s\na: %s", i+1, linePrefix, line)
			}
		}
	}
}

func TestCustomReporterMergesDataWithContext(t *testing.T) {
	context := NewContext(nil)

	exceptions := make(chan *exception, 1)
	context.ExceptionReporter = &channelExceptionReporter{exceptions}

	context.Add("a", 1)
	context.Add("b", 1)

	err := fmt.Errorf("Test")
	context.Report(err, Data{"b": 2})
	exception := <-exceptions

	expectedData := Data{"a": 1, "b": 2}
	if exception.Data["a"] != expectedData["a"] || exception.Data["b"] != expectedData["b"] {
		t.Errorf("Expected exception data to be %v but was %v", expectedData, exception.Data)
	}
}

type exception struct {
	Error error
	Data  Data
}

type channelExceptionReporter struct {
	Channel chan *exception
}

func (c *channelExceptionReporter) Report(err error, data Data) error {
	c.Channel <- &exception{err, data}
	return nil
}
