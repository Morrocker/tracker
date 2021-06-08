package tracker

import (
	"sync/atomic"
)

type Counter interface {
	SetCurrent(n int64)
	Current(n int64) (int64, error)
	RawValue() int64
	Value() string
	UnitsFunc(func(int64) string)
	pointer() *int64
}

type counter struct {
	name      string
	current   int64
	unitsFunc func(int64) string
}

func NewCounter(name string) Counter {
	newCounter := &counter{
		name:    name,
		current: 0,
	}
	return newCounter
}

func (c *counter) SetCurrent(n int64) {
	atomic.CompareAndSwapInt64(&c.current, c.current, n)
}

func (c *counter) Current(n int64) (int64, error) {
	atomic.AddInt64(&c.current, n)
	return c.current, nil
}

func (c *counter) RawValue() int64 {
	return c.current
}
func (c *counter) Value() string {
	if c.unitsFunc == nil {
		return "unitsFunction not set"
	}
	return c.unitsFunc(c.current)
}

func (c *counter) UnitsFunc(f func(int64) string) {
	c.unitsFunc = f
}

func (c *counter) pointer() *int64 {
	return &c.current
}
