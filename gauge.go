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
	ticksLapse  time.Duration
	ticker      *time.Ticker
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

func newGauge(name string, total int64) *gauge {
	var lock sync.Mutex
	newgauge := &gauge{
		name:    name,
		current: 0,
		total:   total,
		show:    defVis,
		mode:    defGaugeMode,
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

func (g *gauge) startAutoMeasure(d time.Duration) error {
	op := "tracker.startAutoMeasure()"
	if g.speedRate == nil {
		return errors.New(op, "speedRate variable not set")
	}
	if err := g.checkTicker(); err != nil {
		g.restartTicker()
	}
	g.ticksLapse = d
	g.ticker = time.NewTicker(d)

	go func() {
		for {
			select {
			case <-g.ticker.C:
				end := g.spdMeasureStart()
				time.Sleep(g.ticksLapse)
				end()
				// time.Sleep(500 * time.Millisecond)
				// r := g.speedRate.AvgRate() // THIS SHOULDNT BE HERE AT ALL
				// r := g.speedRate.Debug()   // THIS SHOULDNT BE HERE AT ALL
				// Logger.Bench("Start: %d | End: %d | Rate: %d | RateTot: %d / SampleSize: %d")
			}
		}
	}()
	return nil
}

func (g *gauge) stopAutoMeasure() error {
	if err := g.checkTicker(); err != nil {
		return errors.Extend("gauge.stopAutoMeasure()", err)
	}
	g.ticker.Stop()
	return nil
}

func (g *gauge) restartTicker() error {
	if g.ticker == nil {
		return errors.New("tracker_group.resetTicker()", "Ticker hasn't been set")
	}
	g.ticker.Reset(g.ticksLapse)
	return nil
}

func (g *gauge) checkTicker() error {
	if g.ticker == nil {
		return errors.New("gauge.checkTicker()", "Ticker hasn't been set")
	}
	return nil
}
