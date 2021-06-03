package tracker

import (
	"sync/atomic"
)

type Gauge interface {
	SetCurrent(n int64)
	Current(n int64) int64
	SetTotal(n int64)
	Total(n int64) int64
	RawValues() (int64, int64)
	Values() (string, string)
	UnitsFunc(func(int64) string)
}

type gauge struct {
	name     string
	current  int64
	total    int64
	unitFunc func(int64) string
}

func NewGauge(name string, total int64) Gauge {
	newGauge := &gauge{
		name:    name,
		current: 0,
		total:   total,
	}
	return newGauge
}

func (g *gauge) SetCurrent(n int64) {
	atomic.CompareAndSwapInt64(&g.current, g.current, n)
}
func (g *gauge) Current(n int64) int64 {
	atomic.AddInt64(&g.current, n)
	return g.current
}
func (g *gauge) SetTotal(n int64) {
	atomic.CompareAndSwapInt64(&g.total, g.total, n)
}
func (g *gauge) Total(n int64) int64 {
	atomic.AddInt64(&g.total, n)
	return g.total
}

func (g *gauge) RawValues() (int64, int64) {
	return g.current, g.total
}
func (g *gauge) Values() (string, string) {
	if g.unitFunc == nil {
		return "unitsFunction not set", "unitsFunction not set"
	}
	return g.unitFunc(g.current), g.unitFunc(g.total)
}

func (g *gauge) UnitsFunc(f func(int64) string) {
	g.unitFunc = f
}
