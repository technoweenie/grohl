package grohl

import (
	"testing"
	"time"
)

func TestLogsCounter(t *testing.T) {
	s, buf := setupLogger()
	s.Counter(1.0, "a", 1, 2)
	if line := logged(buf); line != "metric=a count=1\nmetric=a count=2" {
		t.Errorf("Bad log output: %q", line)
	}
}

func TestLogsTiming(t *testing.T) {
	s, buf := setupLogger()
	dur1, _ := time.ParseDuration("15ms")
	dur2, _ := time.ParseDuration("3s")

	s.Timing(1.0, "a", dur1, dur2)
	if line := logged(buf); line != "metric=a timing=15\nmetric=a timing=3000" {
		t.Errorf("Bad log output: %q", line)
	}
}

func TestLogsGauge(t *testing.T) {
	s, buf := setupLogger()
	s.Gauge(1.0, "a", "1", "2")
	if line := logged(buf); line != "metric=a gauge=1\nmetric=a gauge=2" {
		t.Errorf("Bad log output: %q", line)
	}
}
