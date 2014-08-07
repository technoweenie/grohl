package grohl

import (
	"fmt"
	"testing"
)

func TestLogsError(t *testing.T) {
	reporter, buf := setupLogger(t)
	reporter.Add("a", 1)
	reporter.Add("b", 1)

	err := fmt.Errorf("Test")

	reporter.Report(err, Data{"b": 2, "c": 3, "at": "overwrite me"})
	firstRow := []string{
		"a=1",
		"b=2",
		"c=3",
		"at=exception",
		"class=*errors.errorString",
		"message=Test",
	}

	otherRows := append(firstRow, "~site=")

	for i, line := range buf.Lines() {
		if i == 0 {
			AssertBuiltLine(t, line, firstRow...)
		} else {
			AssertBuiltLine(t, line, otherRows...)
		}
	}
}

func TestCustomReporterMergesDataWithContext(t *testing.T) {
	context := NewContext(nil)

	errors := make(chan *reportedError, 1)
	context.ErrorReporter = &channelErrorReporter{errors}

	context.Add("a", 1)
	context.Add("b", 1)

	err := fmt.Errorf("Test")
	context.Report(err, Data{"b": 2})
	reportedErr := <-errors

	expectedData := Data{"a": 1, "b": 2}
	if reportedErr.Data["a"] != expectedData["a"] || reportedErr.Data["b"] != expectedData["b"] {
		t.Errorf("Expected error data to be %v but was %v", expectedData, reportedErr.Data)
	}
}

type reportedError struct {
	Error error
	Data  Data
}

type channelErrorReporter struct {
	Channel chan *reportedError
}

func (c *channelErrorReporter) Report(err error, data Data) error {
	c.Channel <- &reportedError{err, data}
	return nil
}
