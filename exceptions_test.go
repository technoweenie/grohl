package grohl

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLogsException(t *testing.T) {
	logger, buf := exceptionLoggerWithBuffer()
	logger.AddContext("a", 1)
	logger.AddContext("b", 1)

	err := fmt.Errorf("Test")

	logger.Report(err, LogData{"b": 2, "c": 3, "at": "overwrite me"})
	if line := logged(buf); line != "a=1 b=2 c=3 at=exception class=*errors.errorString message=Test" {
		t.Errorf("Line does not match: %s", line)
	}
}

func exceptionLoggerWithBuffer() (*ExceptionLogger, *bytes.Buffer) {
	logger, buf := loggerWithBuffer()
	return logger.NewExceptionLogger(nil), buf
}
