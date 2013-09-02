package grohl

import (
	"testing"
)

func TestTimer(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddContext("a", "1")
	timer := logger.NewTimer(LogData{"b": "2"})
	timer.Log(LogData{"at": "clobbered", "c": "3"})

	if result := logged(buf); result != "a=1 b=2 at=start\na=1 b=2 at=finish c=3 elapsed=0.000" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestTimerInMS(t *testing.T) {
	logger, buf := loggerWithBuffer()
	logger.AddContext("a", "1")
	timer := logger.NewTimer(LogData{"b": "2"})
	timer.TimeUnit = "ms"
	timer.Log(LogData{"at": "clobbered", "c": "3"})

	expected := "a=1 b=2 at=start\na=1 b=2 at=finish c=3 elapsed=0.001"
	checkedLen := len(expected) - 1
	if result := logged(buf); result[0:checkedLen] != expected[0:checkedLen] {
		t.Errorf("Bad log output: %s", result)
	}
}
