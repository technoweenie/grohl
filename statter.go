package grohl

import (
	"math/rand"
	"time"
)

type Statter struct {
	Logger Logger
}

func NewStatter() *Statter {
	return &Statter{DefaultLogger}
}

func (s *Statter) Counter(sampleRate float32, bucket string, n ...int) {
	if rand.Float32() > sampleRate {
		return
	}

	for _, num := range n {
		s.Logger.Log(LogData{"metric": bucket, "count": num})
	}
}

func (s *Statter) Timing(sampleRate float32, bucket string, d ...time.Duration) {
	if rand.Float32() > sampleRate {
		return
	}

	for _, dur := range d {
		s.Logger.Log(LogData{"metric": bucket, "timing": int64(dur / time.Millisecond)})
	}
}

func (s *Statter) Gauge(sampleRate float32, bucket string, value ...string) {
	if rand.Float32() > sampleRate {
		return
	}

	for _, v := range value {
		s.Logger.Log(LogData{"metric": bucket, "gauge": v})
	}
}
