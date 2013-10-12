package grohl

import (
	"testing"
)

func TestTimerLog(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.Log(Data{"c": "3"})

	if result := logged(buf); result != "a=1 b=2 at=start\na=1 b=2 c=3 elapsed=0.000" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestTimerLogInMS(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.TimeUnit = "ms"
	timer.Log(Data{"c": "3"})

	expected := "a=1 b=2 at=start\na=1 b=2 c=3 elapsed=0.001"
	checkedLen := len(expected) - 3
	if result := logged(buf); result[0:checkedLen] != expected[0:checkedLen] {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestTimerFinish(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.Finish()

	if result := logged(buf); result != "a=1 b=2 at=start\na=1 b=2 at=finish elapsed=0.000" {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestTimerWithStatter(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	statter := NewContext(nil)
	statter.Logger = context.Logger
	timer.SetStatter(statter, 1.0, "bucket")
	timer.Finish()

	expected := "a=1 b=2 at=start\n"
	expected = expected + "a=1 b=2 at=finish elapsed=0.000\n"
	expected = expected + "metric=bucket timing=0"
	if result := logged(buf); result != expected {
		t.Errorf("Bad log output: %s", result)
	}
}

func TestTimerWithNilStatter(t *testing.T) {
	context, buf := setupLogger()
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.SetStatter(nil, 1.0, "bucket")
	timer.Finish()

	expected := "a=1 b=2 at=start\n"
	expected = expected + "a=1 b=2 at=finish elapsed=0.000\n"
	expected = expected + "a=1 b=2 metric=bucket timing=0"
	if result := logged(buf); result != expected {
		t.Errorf("Bad log output: %s", result)
	}
}
