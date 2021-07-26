package tracker

import (
	"sync/atomic"
)

type Gauge interface {
	SetCurrent(n int64)
	Current(n int64) (int64, error)
	SetTotal(n int64)
	Total(n int64) (int64, error)
	RawValues() (int64, int64)
	Values() (string, string)
	UnitsFunc(func(int64) string)
	Reset()
	Pointers() (*int64, *int64)
}

type gauge struct {
	current   int64
	total     int64
	unitsFunc func(int64) string
}

func NewGauge() Gauge {
	newGauge := &gauge{
	}
	return newGauge
}

func (g *gauge) SetCurrent(n int64) {
	atomic.CompareAndSwapInt64(&g.current, g.current, n)
}
func (g *gauge) Current(n int64) (int64, error) {
	atomic.AddInt64(&g.current, n)
	return g.current, nil
}

func (g *gauge) SetTotal(n int64) {
	atomic.CompareAndSwapInt64(&g.total, g.total, n)
}
func (g *gauge) Total(n int64) (int64, error) {
	atomic.AddInt64(&g.total, n)
	return g.total, nil
}

func (g *gauge) RawValues() (int64, int64) {
	return g.current, g.total
}
func (g *gauge) Values() (string, string) {
	if g.unitsFunc == nil {
		return "unitsFunction not set", "unitsFunction not set"
	}
	return g.unitsFunc(g.current), g.unitsFunc(g.total)
}

func (g *gauge) UnitsFunc(f func(int64) string) {
	g.unitsFunc = f
}

func (g *gauge) Reset(){
	atomic.CompareAndSwapInt64(&g.current, g.current, 0)
	atomic.CompareAndSwapInt64(&g.total, g.total, 0)
}

func (g *gauge) Pointers() (*int64, *int64) {
	return &g.current, &g.total
}
