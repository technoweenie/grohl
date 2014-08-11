package grohl

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrapHttpError(t *testing.T) {
	err := errors.New("sup")
	e := NewHttpError(err, 0)
	if e.Message != "sup" {
		t.Errorf("Unexpected error message: %s", e.Message)
	}

	if e.StatusCode != 500 {
		t.Errorf("Unexpected status code: %d", e.StatusCode)
	}

	if e.InnerError != err {
		t.Errorf("Unexpected inner error: %v", e.InnerError)
	}
}

func TestWrapHttpErrorWithStatus(t *testing.T) {
	err := errors.New("sup")
	e := NewHttpError(err, 409)
	if e.Message != "sup" {
		t.Errorf("Unexpected error message: %s", e.Message)
	}

	if e.StatusCode != 409 {
		t.Errorf("Unexpected status code: %d", e.StatusCode)
	}

	if e.InnerError != err {
		t.Errorf("Unexpected inner error: %v", e.InnerError)
	}
}

func TestWrapError(t *testing.T) {
	err := errors.New("sup")
	e := NewError(err)
	if e.Message != "sup" {
		t.Errorf("Unexpected error message: %s", e.Message)
	}

	if e.InnerError != err {
		t.Errorf("Unexpected inner error: %v", e.InnerError)
	}
}

func TestWrapErrorWithMessage(t *testing.T) {
	err := errors.New("sup")
	e := NewErrorf(err, "nuff said, %s", "bub")
	if e.Message != "nuff said, bub" {
		t.Errorf("Unexpected error message: %s", e.Message)
	}

	if e.InnerError != err {
		t.Errorf("Unexpected inner error: %v", e.InnerError)
	}
}

func TestWrapNilError(t *testing.T) {
	e := NewError(nil)
	if e.Message != "" {
		t.Errorf("Expected empty error message: %s", e.Message)
	}

	if e.InnerError != nil {
		t.Errorf("Expected nil inner error: %v", e.InnerError)
	}
}

func TestWrapNilErrorWithMessage(t *testing.T) {
	e := NewErrorf(nil, "nuff said, %s", "bub")
	if e.Message != "nuff said, bub" {
		t.Errorf("Unexpected error message: %s", e.Message)
	}

	if e.InnerError != nil {
		t.Errorf("Expected nil inner error: %v", e.InnerError)
	}
}

func TestLogsWrappedError(t *testing.T) {
	err := errors.New("sup")
	e := NewErrorf(err, "wat")
	e.Add("b", 2)
	e.Add("c", 2)

	reporter, buf := setupLogger(t)
	reporter.Add("a", 1)
	reporter.Add("b", 1)

	reporter.Report(e, Data{"c": 3, "d": 4, "at": "overwrite"})
	firstRow := []string{
		"a=1",
		"b=2",
		"c=3",
		"d=4",
		"at=exception",
		"class=*grohl.Err",
		"message=sup",
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
