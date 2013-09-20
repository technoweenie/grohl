package grohl

import (
	"testing"
)

func TestTimer(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.Log(Data{"at": "clobbered", "c": "3"})

	if result := logged(buf); result != "a=1 b=2 at=start\na=1 b=2 at=finish c=3 elapsed=0.000" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestTimerInMS(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.TimeUnit = "ms"
	timer.Log(Data{"at": "clobbered", "c": "3"})

	expected := "a=1 b=2 at=start\na=1 b=2 at=finish c=3 elapsed=0.001"
	checkedLen := len(expected) - 3
	if result := logged(buf); result[0:checkedLen] != expected[0:checkedLen] {
		t.Errorf("Bad log output: %s", result)
	}
}
