package tracker

import (
	"time"

	"github.com/morrocker/benchmark"
)

type Speed interface {
	SampleSize(uint) uint
	Reset()
	StartMeasure() func()
	StartAutoMeasure(time.Duration)
	StopAutoMeasure()
	UnitsFunc(func(int64) string)
	RawRate() int64
	Rate() string
}

type speed struct {
	target    *int64
	ticker    *time.Ticker
	rate      benchmark.SingleRate
	unitsFunc func(int64) string
}

func NewSpeed(g *int64, n uint) Speed {
	newSpeed := &speed{
		target: g,
		rate:   benchmark.NewSingleRate(n),
	}
	return newSpeed
}

func (s *speed) SampleSize(n uint) uint {
	return s.rate.SampleSize(n)
}

func (s *speed) Reset() {
	s.rate.Reset()
}

func (s *speed) StartMeasure() func() {
	end := s.rate.MeasureStart(*s.target)
	return func() {
		end(*s.target)
	}
}

func (s *speed) StartAutoMeasure(d time.Duration) {
	if s.ticker != nil {
		s.ticker.Reset(d)
		return
	}
	s.ticker = time.NewTicker(d)

	go func() {
		for range s.ticker.C {
			end := s.StartMeasure()
			time.Sleep(d)
			end()
		}
	}()
}

func (s *speed) StopAutoMeasure() {
	s.ticker.Stop()
}

func (s *speed) RawRate() int64 {
	return s.rate.AvgRate()
}

func (s *speed) Rate() string {
	return s.unitsFunc(s.rate.AvgRate())
}

func (s *speed) UnitsFunc(fn func(int64) string) {
	s.unitsFunc = fn
}
