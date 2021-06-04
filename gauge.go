package tracker

import (
	"sync/atomic"

	"github.com/morrocker/errors"
)

type Gauge interface {
	SetCurrent(n uint64)
	Current(n int64) (uint64, error)
	SetTotal(n uint64)
	Total(n int64) (uint64, error)
	RawValues() (uint64, uint64)
	Values() (string, string)
	UnitsFunc(func(uint64) string)
}

type gauge struct {
	name      string
	current   uint64
	total     uint64
	unitsFunc func(uint64) string
}

func NewGauge(name string, total uint64) Gauge {
	newGauge := &gauge{
		name:    name,
		current: 0,
		total:   total,
	}
	return newGauge
}

func (g *gauge) SetCurrent(n uint64) {
	atomic.CompareAndSwapUint64(&g.current, g.current, n)
}
func (g *gauge) Current(n int64) (uint64, error) {
	if n < 0 {
		if (int64(g.current) + n) < 0 {
			return 0, errors.New("tracker.gauge.Current()", "Underflow error. Current substraction value is negative")
		}
		atomic.AddUint64(&g.current, -uint64(n))
	} else {
		atomic.AddUint64(&g.current, uint64(n))
	}
	return g.current, nil
}

func (g *gauge) SetTotal(n uint64) {
	atomic.CompareAndSwapUint64(&g.total, g.total, n)
}
func (g *gauge) Total(n int64) (uint64, error) {
	if n < 0 {
		if (int64(g.total) + n) < 0 {
			return 0, errors.New("tracker.gauge.Current()", "Underflow error. Current substraction value is negative")
		}
		atomic.AddUint64(&g.total, -uint64(n))
	} else {
		atomic.AddUint64(&g.total, uint64(n))
	}
	return g.total, nil
}

func (g *gauge) RawValues() (uint64, uint64) {
	return g.current, g.total
}
func (g *gauge) Values() (string, string) {
	if g.unitsFunc == nil {
		return "unitsFunction not set", "unitsFunction not set"
	}
	return g.unitsFunc(g.current), g.unitsFunc(g.total)
}

func (g *gauge) UnitsFunc(f func(uint64) string) {
	g.unitsFunc = f
}
