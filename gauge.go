package tracker

import (
	"fmt"
	"sync"
	"time"

	"github.com/morrocker/benchmark"
	"github.com/morrocker/errors"
)

type gauge struct {
	name        string
	current     int64
	total       int64
	show        bool
	mode        string
	order       int
	autoMeasure bool
	lock        *sync.Mutex
	speedRate   *benchmark.SRate
	unitFunc    func(int64) string
}

var modes = []string{
	"countdown",
	"division",
	"percentageDone",
	"percentageRemaining",
	"countPercentage",
	"divisionPercentage",
}

func newGauge(name string, total int64, order int) *gauge {
	var lock sync.Mutex
	newgauge := &gauge{
		name:    name,
		current: 0,
		total:   total,
		show:    defVis,
		mode:    defGaugeMode,
		order:   order,
		lock:    &lock,
	}
	return newgauge
}

func (g *gauge) setCurrent(n int64) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.current = n
}
func (g *gauge) changeCurrent(n int64) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.current += n
}
func (g *gauge) setTotal(n int64) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.total = n
}
func (g *gauge) changeTotal(n int64) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.total += n
}
func (g *gauge) setMode(mode string) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.mode = mode
}
func (g *gauge) getMode() string {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.mode
}
func (g *gauge) setOrder(n int) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.order = n
}
func (g *gauge) getOrder() int {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.order
}
func (g *gauge) getRawValues() (int64, int64) {
	g.lock.Lock()
	defer g.lock.Unlock()
	return g.current, g.total
}
func (g *gauge) getValues() (string, string, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.unitFunc == nil {
		return "", "", errors.New("gauge.getValues()", "unitsFunction not set")
	}
	return g.unitFunc(g.current), g.unitFunc(g.total), nil
}

func (g *gauge) spdMeasureStart() func() {
	end := g.speedRate.MeasureStart(g.current)
	return func() {
		end(g.current)
	}
}

func (g *gauge) getRawRate() int64 {
	return g.speedRate.AvgRate()
}

func (g *gauge) setUnitsFunc(f func(int64) string) {
	g.unitFunc = f
}

func (g *gauge) getRate() string {
	if g.unitFunc != nil {
		return g.unitFunc(g.speedRate.AvgRate())
	}
	return fmt.Sprintf("%d UntypedUnit", g.speedRate.AvgRate())
}

func (g *gauge) getETA() string {
	if x := g.speedRate.AvgRate(); x != 0 {
		eta := (g.total - g.current) / x
		return time.Duration(eta * 1000000000).String()
	}
	return "not available"
}

func (g *gauge) initSpdRate(n int) {
	spd := benchmark.NewSRate(n)
	g.speedRate = spd
}

func (g *gauge) startAutoMeasure(tick int) error {
	op := "tracker.startAutoMeasure()"
	if g.speedRate == nil {
		return errors.New(op, "speedRate variable not set")
	}
	if g.autoMeasure {
		return nil
	}
	g.autoMeasure = true
	go func() {
		for {
			measureEnd := g.spdMeasureStart()
			time.Sleep(time.Duration(tick) * time.Second)
			measureEnd()
			if !g.autoMeasure {
				break
			}
		}
	}()
	return nil
}

func (g *gauge) stopAutoMeasure() error {
	op := "tracker.startAutoMeasure()"
	if g.speedRate == nil {
		return errors.New(op, "speedRate variable not set")
	} else if !g.autoMeasure {
		return nil
	}
	g.autoMeasure = false
	return nil
}

func (g *gauge) print() (format string) {
	g.lock.Lock()
	defer g.lock.Unlock()

	switch g.mode {
	case "countdown":
		diff := (g.total - g.current)
		format = fmt.Sprintf("%s: %d remaining", g.name, diff)
	case "division":
		if g.unitFunc != nil {
			format = fmt.Sprintf("%s: %s/%s", g.name, g.unitFunc(g.current), g.unitFunc(g.total))
			return
		}
		format = fmt.Sprintf("%s: %d/%d", g.name, g.current, g.total)
	case "percentageDone":
		per := g.current * 100 / g.total
		format = fmt.Sprintf("%s: %d %% done", g.name, per)
	case "percentageRemaining":
		per := (g.total - g.current) * 100 / g.total
		format = fmt.Sprintf("%s: %d %% remaining", g.name, per)
	case "countPercentage":
		per := g.current * 100 / g.total
		format = fmt.Sprintf("%s: %d %% remaining", g.name, per)
	case "divisionPercentage":
		per := g.current * 100 / g.total
		if g.unitFunc != nil {
			format = fmt.Sprintf("%s: %s/%s ( %d %% remaining )", g.name, g.unitFunc(g.current), g.unitFunc(g.total), per)
			return
		}
		format = fmt.Sprintf("%s: %d/%d ( %d %% remaining )", g.name, g.current, g.total, per)
	}
	return
}
