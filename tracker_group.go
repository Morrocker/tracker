package tracker

import (
	"time"

	"github.com/morrocker/errors"
	"github.com/morrocker/log"
)

type Group interface {
	AddGauge(string)
	AddCounter(string)
	AddSpeed(string, *int64, uint)
	Counter(string) Counter
	Gauge(string) Gauge
	Speed(string) Speed
	PrintFunc(func())
	Print()
	StartAutoPrint(time.Duration)
	StopAutoPrint()
	RestartAutoPrint()
}

type group struct {
	gauges      map[string]Gauge
	counters    map[string]Counter
	speeds      map[string]Speed
	prntFunc    func()
	ticker      *time.Ticker
	tickerCicle time.Duration
}

func NewGroup() Group {
	newGroup := &group{
		gauges:   make(map[string]Gauge),
		counters: make(map[string]Counter),
		speeds:   make(map[string]Speed),
	}
	return newGroup
}

func (g *group) AddGauge(name string) {
	g.gauges[name] = NewGauge(name)
}

func (g *group) Gauge(s string) Gauge {
	ret, ok := g.gauges[s]
	if !ok {
		log.Errorln(errors.New("tracker.group.Gauge()", "gauge not found"))
		return nil
	}
	return ret
}

func (g *group) AddCounter(name string) {
	g.counters[name] = NewCounter(name)
}

func (g *group) Counter(s string) Counter {
	ret, ok := g.counters[s]
	if !ok {
		log.Errorln(errors.New("tracker.group.Counter()", "counter not found"))
		return nil
	}
	return ret
}

func (g *group) AddSpeed(name string, ptr *int64, sampleSize uint) {
	g.speeds[name] = NewSpeed(ptr, sampleSize)
}

func (g *group) Speed(s string) Speed {
	ret, ok := g.speeds[s]
	if !ok {
		log.Errorln(errors.New("tracker.group.Speed()", "speed not found"))
		return nil
	}
	return ret
}

func (g *group) PrintFunc(fn func()) {
	g.prntFunc = fn
}

func (g *group) Print() {
	g.prntFunc()
}

func (g *group) StartAutoPrint(t time.Duration) {
	if g.ticker != nil {
		g.ticker.Reset(t)
		return
	}
	g.ticker = time.NewTicker(t)

	go func() {
		for range g.ticker.C {
			g.Print()
		}
	}()
}

func (g *group) StopAutoPrint() {
	g.ticker.Stop()
}
func (g *group) RestartAutoPrint() {
	g.ticker.Reset(g.tickerCicle)
}
