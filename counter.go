package tracker

import (
	"sync/atomic"
)

type Counter interface {
	SetCurrent(n int64)
	Current(n int64) int64
	RawValue() int64
	Value() string
	UnitsFunc(func(int64) string)
}

type counter struct {
	name     string
	current  int64
	unitFunc func(int64) string
}

func newCounter(name string, total int64) Counter {
	newCounter := &counter{
		name:    name,
		current: 0,
	}
	return newCounter
}

func (g *counter) SetCurrent(n int64) {
	atomic.CompareAndSwapInt64(&g.current, g.current, n)
}
func (g *counter) Current(n int64) int64 {
	atomic.AddInt64(&g.current, n)
	return g.current
}

func (g *counter) RawValue() int64 {
	return g.current
}
func (g *counter) Value() string {
	if g.unitFunc == nil {
		return "unitsFunction not set"
	}
	return g.unitFunc(g.current)
}

func (g *counter) UnitsFunc(f func(int64) string) {
	g.unitFunc = f
}
