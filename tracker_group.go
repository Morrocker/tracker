package tracker

import (
	"fmt"
	"time"

	"github.com/morrocker/errors"
)

type trackerGroup struct {
	trackers   map[string]tracker
	order      []string
	etaTracker string
	delayPrint bool
	ticksLapse time.Duration
	format     format
	autoPrint  bool
	ticker     *time.Ticker
	printFunc  func() string
}

type format struct {
	showStatus bool
	separator  string
	status     string
	lineMode   string
}

// NewGroup creates a new tracker group within a SuperTracker.
func newGroup(name string) *trackerGroup {
	group := &trackerGroup{
		trackers: make(map[string]tracker),
		format: format{
			lineMode:  defLineMode,
			separator: defSeparator,
		},
		ticksLapse: defLapse,
	}
	return group
}

// AddGauge call AddGaugeOn using the default group
func (g *trackerGroup) addGauge(trackerName, printName string, total interface{}) error {
	op := "tracker_group.addGauge()"
	if _, err := g.findTracker(trackerName); err == nil {
		return errors.New(op, fmt.Sprintf("tracker name %s already taken", trackerName))
	}
	total64, err := getInt64(total)
	if err != nil {
		return errors.New(op, err)
	}
	g.trackers[trackerName] = newGauge(trackerName, total64, 0) // FIX ORDER
	return nil
}

func (g *trackerGroup) findTracker(t string) (tracker, error) {
	tracker := &gauge{}
	for name, tracker := range g.trackers {
		if t == name {
			return tracker, nil
		}
	}
	err := errors.New("tracker.findTracker()", "Did not find tracker "+t)
	return tracker, err
}

func (g *trackerGroup) setLineMode(mode string) error {
	switch mode {
	case "singleline":
		g.format.lineMode = "singleline"
	case "multiline":
		g.format.lineMode = "multiline"
	default:
		return errors.New("tracker_group.SetLineMode()", "Set mode %s is not a valid mode")
	}
	return nil
}

func (g *trackerGroup) setEtaTracker(tracker string) error {
	op := "tracker_group.SetEtaTracker()"
	if _, err := g.findTracker(tracker); err != nil {
		return errors.New(op, "Did not find tracker "+tracker)
	}
	g.etaTracker = tracker
	return nil
}

func (g *trackerGroup) changeCurr(tracker string, value interface{}) error {
	op := "tracker_group.SetEtaTracker()"
	Tracker, err := g.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}
	val64, err := getInt64(value)
	if err != nil {
		return errors.Extend(op, err)
	}
	Tracker.changeCurrent(val64)
	return nil
}

func (g *trackerGroup) setPrintFunc(f func() string) {
	g.printFunc = f
}

func (g *trackerGroup) print() {
	g.printFunc()
}

func (g *trackerGroup) startAutoPrint(d time.Duration) {
	if err := g.checkTicker(); err == nil {
		g.restartTicker()
	}
	g.ticksLapse = d
	g.ticker = time.NewTicker(d)

	go func() {
		select {
		case <-g.ticker.C:
			g.print()
		}
	}()
}

func (g *trackerGroup) stopAutoPrint() error {
	if err := g.checkTicker(); err != nil {
		return errors.Extend("tracker_group.resetTicker()", err)
	}
	g.ticker.Stop()
	return nil
}

func (g *trackerGroup) restartTicker() error {
	if g.ticker == nil {
		return errors.New("tracker_group.resetTicker()", "Ticker hasn't been set")
	}
	g.ticker.Reset(g.ticksLapse)
	return nil
}

func (g *trackerGroup) checkTicker() error {
	if g.ticker == nil {
		return errors.New("tracker_group.checkTicker()", "Ticker hasn't been set")
	}
	return nil
}
