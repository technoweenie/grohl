package grohl

import (
	"strconv"
	"time"
)

type timer struct {
	Started time.Time
	*IoLogger
}

// A timer is a special logger that tracks the duration spent since its creation.
func (l *IoLogger) NewTimer(context map[string]interface{}) *timer {
	ctx := l.NewContext(context)
	ctx.Log(LogData{"at": "start"})
	return &timer{time.Now(), ctx}
}

// Writes a final log message with the elapsed time shown.
func (t *timer) Finish() {
	t.Log(nil)
}

// Writes a log message with extra data or the elapsed time shown.  Pass nil or
// use Finish() if there is no extra data.
func (t *timer) Log(data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	data["at"] = "finish"
	data["elapsed"] = t.ElapsedString()
	t.IoLogger.Log(data)
}

func (t *timer) Elapsed() time.Duration {
	return time.Since(t.Started)
}

func (t *timer) ElapsedString() string {
	dur := t.Elapsed()
	return strconv.FormatFloat(t.durationUnit(dur), durationFormat, 3, 64)
}

func (t *timer) durationUnit(dur time.Duration) float64 {
	sec := dur.Seconds()
	if t.TimeUnit == "ms" {
		return sec * 1000
	}
	return sec
}

var durationFormat = []byte("f")[0]
