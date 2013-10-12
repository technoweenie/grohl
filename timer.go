package grohl

import (
	"time"
)

type Timer struct {
	Started        time.Time
	TimeUnit       string
	context        *Context
	statter        Statter
	statSampleRate float32
	statBucket     string
}

// A timer tracks the duration spent since its creation.
func (c *Context) Timer(data Data) *Timer {
	context := c.New(data)
	context.Log(Data{"at": "start"})
	timer := &Timer{Started: time.Now(), TimeUnit: context.TimeUnit, context: context}
	if c.statter != nil {
		timer.SetStatter(c.statter, c.statSampleRate, c.statBucket)
	}
	return timer
}

// Writes a final log message with the elapsed time shown.
func (t *Timer) Finish() {
	dur := t.Elapsed()

	t.Log(Data{"at": "finish", "elapsed": t.durationUnit(dur)})

	if t.statter != nil {
		t.statter.Timing(t.statSampleRate, t.statBucket, dur)
	}
}

// Writes a log message with extra data or the elapsed time shown.  Pass nil or
// use Finish() if there is no extra data.
func (t *Timer) Log(data Data) error {
	if data == nil {
		data = make(Data)
	}

	if _, ok := data["elapsed"]; !ok {
		data["elapsed"] = t.durationUnit(t.Elapsed())
	}

	return t.context.Log(data)
}

func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.Started)
}

func (t *Timer) SetStatter(statter Statter, sampleRate float32, bucket string) {
	if statter == nil {
		t.statter = CurrentStatter
	} else {
		t.statter = statter
	}
	t.statSampleRate = sampleRate
	t.statBucket = bucket
}

func (t *Timer) durationUnit(dur time.Duration) float64 {
	sec := dur.Seconds()
	if t.TimeUnit == "ms" {
		return sec * 1000
	}
	return sec
}
