package grohl

import (
	"time"
)

type Timer struct {
	Started  time.Time
	TimeUnit string
	context  *Context
	*_statter
}

// A timer tracks the duration spent since its creation.
func (c *Context) Timer(data Data) *Timer {
	context := c.New(data)
	context.Log(Data{"at": "start"})
	return &Timer{
		Started:  time.Now(),
		TimeUnit: context.TimeUnit,
		context:  context,
		_statter: c._statter,
	}
}

// Writes a final log message with the elapsed time shown.
func (t *Timer) Finish() {
	dur := t.Elapsed()

	t.Log(Data{"at": "finish", "elapsed": t.durationUnit(dur)})
	t._statter.Timing(dur)
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

func (t *Timer) durationUnit(dur time.Duration) float64 {
	sec := dur.Seconds()
	if t.TimeUnit == "ms" {
		return sec * 1000
	}
	return sec
}
