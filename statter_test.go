package grohl

import (
	"testing"
	"time"
)

func TestLogsCounter(t *testing.T) {
	s, buf := setupLogger(t)
	s.Counter(1.0, "a", 1, 2)
	buf.AssertLogged("metric=a count=1\nmetric=a count=2")
}

func TestLogsTiming(t *testing.T) {
	s, buf := setupLogger(t)
	dur1, _ := time.ParseDuration("15ms")
	dur2, _ := time.ParseDuration("3s")

	s.Timing(1.0, "a", dur1, dur2)
	buf.AssertLogged("metric=a timing=15\nmetric=a timing=3000")
}

func TestLogsGauge(t *testing.T) {
	s, buf := setupLogger(t)
	s.Gauge(1.0, "a", "1", "2")
	buf.AssertLogged("metric=a gauge=1\nmetric=a gauge=2")
}
