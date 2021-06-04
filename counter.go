package tracker

import (
	"sync/atomic"

	"github.com/morrocker/errors"
)

type Counter interface {
	SetCurrent(n uint64)
	Current(n int64) (uint64, error)
	RawValue() uint64
	Value() string
	UnitsFunc(func(uint64) string)
}

type counter struct {
	name      string
	current   uint64
	unitsFunc func(uint64) string
}

func NewCounter(name string) Counter {
	newCounter := &counter{
		name:    name,
		current: 0,
	}
	return newCounter
}

func (g *counter) SetCurrent(n uint64) {
	atomic.CompareAndSwapUint64(&g.current, g.current, n)
}

func (g *counter) Current(n int64) (uint64, error) {
	if n < 0 {
		if (int64(g.current) + n) < 0 {
			return 0, errors.New("tracker.counter.Current()", "Underflow error. Current substraction value is negative")
		}
		atomic.AddUint64(&g.current, -uint64(n))
	} else {
		atomic.AddUint64(&g.current, uint64(n))
	}
	return g.current, nil
}

func (g *counter) RawValue() uint64 {
	return g.current
}
func (g *counter) Value() string {
	if g.unitsFunc == nil {
		return "unitsFunction not set"
	}
	return g.unitsFunc(g.current)
}

func (g *counter) UnitsFunc(f func(uint64) string) {
	g.unitsFunc = f
}
