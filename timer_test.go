package grohl

import (
	"testing"
)

func TestTimerLog(t *testing.T) {
	context, buf := setupLogger(t)
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.Add("c", "3")
	timer.Log(Data{"d": "4"})

	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("a=1", "b=2", "c=3", "d=4", "elapsed=0.000")
	buf.AssertEOF()
}

func TestTimerLogInMS(t *testing.T) {
	context, buf := setupLogger(t)
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.TimeUnit = "ms"
	timer.Log(Data{"c": "3"})

	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("a=1", "b=2", "c=3", "~elapsed=0.00")
	buf.AssertEOF()
}

func TestTimerFinish(t *testing.T) {
	context, buf := setupLogger(t)
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.Add("c", "3")
	timer.Finish()

	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("a=1", "b=2", "c=3", "at=finish", "elapsed=0.000")
	buf.AssertEOF()
}

func TestTimerWithCurrentStatter(t *testing.T) {
	context, buf := setupLogger(t)
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	timer.StatterBucketSuffix("bucket")

	oldStatter := CurrentStatter
	CurrentStatter = context
	timer.Finish()
	CurrentStatter = oldStatter

	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("a=1", "metric=bucket", "timing=0")
	buf.AssertLine("a=1", "b=2", "at=finish", "elapsed=0.000")
}

func TestTimerWithStatter(t *testing.T) {
	context, buf := setupLogger(t)
	context.Add("a", "1")
	timer := context.Timer(Data{"b": "2"})
	statter := NewContext(nil)
	statter.Logger = context.Logger
	timer.SetStatter(statter, 1.0, "bucket")
	timer.Finish()

	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("metric=bucket", "timing=0")
	buf.AssertLine("a=1", "b=2", "at=finish", "elapsed=0.000")
	buf.AssertEOF()
}

func TestTimerWithContextStatter(t *testing.T) {
	context, buf := setupLogger(t)
	context.Add("a", "1")
	context.SetStatter(context, 1.0, "bucket")
	timer := context.Timer(Data{"b": "2"})
	timer.StatterBucket = "bucket2"
	timer.Finish()

	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("a=1", "metric=bucket2", "timing=0")
	buf.AssertLine("a=1", "b=2", "at=finish", "elapsed=0.000")
	buf.AssertEOF()

	if context.StatterBucket == "bucket2" {
		t.Errorf("Context's stat bucket was changed")
	}
}

func TestTimerWithNilStatter(t *testing.T) {
	oldlogger := CurrentContext.Logger

	context, buf := setupLogger(t)
	context.Add("a", "1")
	CurrentContext.Logger = context.Logger
	timer := context.Timer(Data{"b": "2"})
	timer.SetStatter(nil, 1.0, "bucket")
	timer.Finish()

	CurrentContext.Logger = oldlogger
	buf.AssertLine("a=1", "b=2", "at=start")
	buf.AssertLine("metric=bucket", "timing=0")
	buf.AssertLine("a=1", "b=2", "at=finish", "elapsed=0.000")
	buf.AssertEOF()
}
