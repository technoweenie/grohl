package grohl

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestLogsException(t *testing.T) {
	logger, buf := exceptionLoggerWithBuffer()
	logger.AddContext("a", 1)
	logger.AddContext("b", 1)

	err := fmt.Errorf("Test")
	errid := int64(reflect.ValueOf(err).Pointer())

	logger.Report(err, LogData{"b": 2, "c": 3, "at": "overwrite me"})
	expected := fmt.Sprintf("a=1 b=2 c=3 at=exception class=*errors.errorString message=Test exception_id=%d", errid)
	if line := logged(buf); line != expected {
		t.Errorf("Line does not match:\ne: %s\na: %s", expected, line)
	}
}

func exceptionLoggerWithBuffer() (*ExceptionLogger, *bytes.Buffer) {
	logger, buf := loggerWithBuffer()
	return logger.NewExceptionLogger(nil), buf
}
