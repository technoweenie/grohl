package grohl

import (
	"time"
)

type Timer struct {
	Started time.Time
	*IoLogger
}

// A timer is a special logger that tracks the duration spent since its creation.
func (l *IoLogger) NewTimer(context map[string]interface{}) *Timer {
	ctx := l.NewContext(context)
	ctx.Log(LogData{"at": "start"})
	return &Timer{time.Now(), ctx}
}

// Writes a final log message with the elapsed time shown.
func (t *Timer) Finish() {
	t.Log(nil)
}

// Writes a log message with extra data or the elapsed time shown.  Pass nil or
// use Finish() if there is no extra data.
func (t *Timer) Log(data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	data["at"] = "finish"
	data["elapsed"] = t.durationUnit(t.Elapsed())
	t.IoLogger.Log(data)
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
