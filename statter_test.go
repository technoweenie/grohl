package grohl

import (
	"testing"
	"time"
)

func TestLogsCounter(t *testing.T) {
	s, buf := setupLogger(t)
	s.Counter(1.0, "a", 1, 2)
	buf.AssertLine("metric=a", "count=1")
	buf.AssertLine("metric=a", "count=2")
	buf.AssertEOF()
}

func TestLogsTiming(t *testing.T) {
	s, buf := setupLogger(t)
	dur1, _ := time.ParseDuration("15ms")
	dur2, _ := time.ParseDuration("3s")

	s.Timing(1.0, "a", dur1, dur2)
	buf.AssertLine("metric=a", "timing=15")
	buf.AssertLine("metric=a", "timing=3000")
	buf.AssertEOF()
}

func TestLogsGauge(t *testing.T) {
	s, buf := setupLogger(t)
	s.Gauge(1.0, "a", "1", "2")
	buf.AssertLine("metric=a", "gauge=1")
	buf.AssertLine("metric=a", "gauge=2")
	buf.AssertEOF()
}

var suffixTests = []string{"abc", "abc."}

func TestSetsBucketSuffix(t *testing.T) {
	s, _ := setupLogger(t)
	for _, prefix := range suffixTests {
		s.StatterBucket = prefix
		s.StatterBucketSuffix("def")
		if s.StatterBucket != "abc.def" {
			t.Errorf("bucket is wrong after prefix %s: %s", prefix, s.StatterBucket)
		}
	}
}
