package grohl

import (
	"fmt"
	"strings"
	"testing"
)

func TestLogsException(t *testing.T) {
	reporter, buf := setupLogger()
	reporter.Add("a", 1)
	reporter.Add("b", 1)

	err := fmt.Errorf("Test")

	reporter.Report(err, Data{"b": 2, "c": 3, "at": "overwrite me"})
	expected := fmt.Sprintf("a=1 b=2 c=3 at=exception class=*errors.errorString message=Test exception_id=%s", ErrorId(err))
	linePrefix := expected + " site="

	for i, line := range strings.Split(logged(buf), "\n") {
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
