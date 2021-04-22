package tracker

import (
	"fmt"
	"time"

	"github.com/morrocker/errors"
	"github.com/morrocker/log"
)

const (
	defLapse = 6000 * time.Millisecond
	defGroup = "default"
	defVis   = true
)

// SuperTracker base structure that houses all tracks to monitor task progress
type SuperTracker struct {
	trackerGroups map[string]*trackerGroup
}

type tracker interface {
	setCurrent(int64)
	setTotal(int64)
	changeCurrent(int64)
	changeTotal(int64)
	changeAndReturn(int64, int64) (int64, int64)
	getRawValues() (int64, int64)
	getValues() (string, string, error)
	spdMeasureStart() func()
	setUnitsFunc(func(int64) string)
	getRawRate() int64
	getRate() string
	getETA() string
	initSpdRate(int)
	startAutoMeasure(time.Duration) error
	stopAutoMeasure() error
}

// New creates a new SuperTracker. The 'default' group is created with it.
func New() *SuperTracker {
	tracker := &SuperTracker{
		trackerGroups: make(map[string]*trackerGroup),
	}
	tracker.trackerGroups[defGroup] = newGroup(defGroup)
	return tracker
}

// AddGauge creates a new gauge on the 'default' group
func (t *SuperTracker) AddGauge(tracker, printName string, total interface{}, group ...string) error {
	op := "tracker.AddGauge()"
	tGroup, err := t.findGroup(group[0])
	if err != nil {
		return errors.Extend(op, err)
	}
	if err := tGroup.addGauge(tracker, printName, total); err != nil {
		return errors.Extend(op, err)
	}
	return nil
}

// ChangeCurr changes a trackers current value by the set amount
func (t *SuperTracker) ChangeCurr(tracker string, value interface{}) (err error) {
	op := "tracker.ChangeCurr()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}
	val64, err := getInt64(value)
	if err != nil {
		return errors.Extend(op, err)
	}
	trckr.changeCurrent(val64)
	return
}

// SetCurr sets a trackers current value to the one given
func (t *SuperTracker) Curr(tracker string, value interface{}) error {
	op := "tracker.IncreaseCurr()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}
	val64, err := getInt64(value)
	if err != nil {
		return errors.Extend(op, err)
	}
	trckr.setCurrent(val64)
	return nil
}

func (t *SuperTracker) ChangeAndReturn(tracker string, values ...interface{}) error {
	op := "tracker.ChangeAndReturn()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}

	curr, err := getInt64(values[0])
	if err != nil {
		return errors.Extend(op, err)
	}
	tot, err := getInt64(values[1])
	if err != nil {
		return errors.Extend(op, err)
	}
	trckr.changeAndReturn(curr, tot)
	return nil
}

// Reset sets a trackers current and total value to 0
func (t *SuperTracker) Reset(tracker string) error {
	op := "tracker.ResetCurr()"
	if err := t.Curr(tracker, 0); err != nil {
		return errors.New(op, err)
	}
	if err := t.Total(tracker, 0); err != nil {
		return errors.New(op, err)
	}
	return nil
}

// ChangeTotal changes a trackers total value by the given amount
func (t *SuperTracker) ChangeTotal(tracker string, value interface{}) error {
	op := "tracker.ChangeTotal()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}
	val64, err := getInt64(value)
	if err != nil {
		return errors.Extend(op, err)
	}
	trckr.changeTotal(val64)
	return nil
}

// Total sets a trackers value to the one given
func (t *SuperTracker) Total(tracker string, value interface{}) error {
	op := "tracker.SetTotal()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}
	val64, err := getInt64(value)
	if err != nil {
		return errors.Extend(op, err)
	}
	trckr.setTotal(val64)
	return nil
}

// Values returns the current and total values from the given tracker
func (t *SuperTracker) Values(tracker string) (current string, total string, err error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		err = errors.Extend("tracker.GetValues()", err)
		return
	}
	return trckr.getValues()
}

// RawValues returns the current and total values from the given tracker
func (t *SuperTracker) RawValues(tracker string) (current int64, total int64, err error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		err = errors.Extend("tracker.GetValues()", err)
		return
	}
	current, total = trckr.getRawValues()
	return
}

// PrintFunc executes the print function for the 'default' group
func (t *SuperTracker) PrintFunc(f func(), group ...string) error {
	tGroup, err := t.findGroup(group...)
	if err != nil {
		log.Errorln(errors.Extend("tracker.SetGroupPrintFunc()", err))
	}
	tGroup.printFunc(f)
	return nil
}

// PrintGroup prints progress of all visible trackers
func (t *SuperTracker) Print(group string) {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.PrintGroup()", err))
	}
	tGroup.print()
}

// Status changes the 'default' group's status
func (t *SuperTracker) Status(task string, group ...string) error {
	tGroup, err := t.findGroup(group[0])
	if err != nil {
		return errors.Extend("tracker.PrintGroup()", err)
	}
	tGroup.setStatus(task)
	return nil
}

// StartMeasure starts a parameter measuring, returning a function that ends the measure
func (t *SuperTracker) StartMeasure(tracker string) (func(), error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return nil, errors.Extend("tracker.StartMeasure()", err)
	}
	endMeasure := trckr.spdMeasureStart()
	return endMeasure, nil
}

// ProgressRate returns the measures average rate without applying the units modifier function
func (t *SuperTracker) ProgressRate(tracker string) (string, error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return "", errors.Extend("tracker.GetProgressRate()", err)
	}
	return trckr.getRate(), nil
}

// TrueProgressRate returns the measure average rate applying the units modifier function
func (t *SuperTracker) TrueProgressRate(tracker string) (string, error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		err = errors.Extend("tracker.GetTrueProgressRate()", err)
		return "", err
	}
	return trckr.getRate(), nil
}

// UnitsFunc sets the given function to be used to modify the tracked units into a human
// readable form
func (t *SuperTracker) UnitsFunc(tracker string, f func(int64) string) error {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend("tracker.SetProgressFunction()", err)
	}
	trckr.setUnitsFunc(f)
	return nil
}

// InitSpdRate TODO STILL THINKING IF having to initialize this separately is a good idea.
func (t *SuperTracker) InitSpdRate(tracker string, n int) error {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend("tracker.IncreaseCurr()", err)
	}
	trckr.initSpdRate(n)
	return nil
}

// StartAutoMeasure starts the automeasuring process. The process will take measures for each
// time.Duration given for the given tracker. TODO automeasure still uses old way to count time
func (t *SuperTracker) StartAutoMeasure(tracker string, tick time.Duration) error {
	op := "tracker.StartAutoMeasure()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}

	if err := trckr.startAutoMeasure(tick); err != nil {
		return errors.Extend(op, err)
	}
	return nil

}

// StopAutoMeasure stops the automeasuring process for the given tracker
func (t *SuperTracker) StopAutoMeasure(tracker string) error {
	op := "tracker.StartAutoMeasure()"
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend(op, err)
	}

	if err := trckr.stopAutoMeasure(); err != nil {
		return errors.Extend(op, err)
	}
	return nil

}

// StartGroupAutoPrint starts the autoprint process for the given group
func (t *SuperTracker) StartAutoPrint(d time.Duration, group ...string) error {
	tGroup, err := t.findGroup(group...)
	if err != nil {
		log.Errorln(errors.Extend("tracker.StartAutoPrintGroup()", err))
	}
	tGroup.startAutoPrint(d)
	return nil
}

// StopGroupAutoPrint stops the autoprint process for the given group
func (t *SuperTracker) StopAutoPrint(group ...string) error {
	tGroup, err := t.findGroup(group...)
	if err != nil {
		log.Errorln(errors.Extend("tracker.StopAutoPrintGroup()", err))
	}
	tGroup.stopAutoPrint()
	return nil
}

// RestartGroupAutoPrint delays the auto-printing for the given group. Also restart's it
// if it was stopped
func (t *SuperTracker) RestartGroupAutoPrint(group ...string) error {
	tGroup, err := t.findGroup(group...)
	if err != nil {
		log.Errorln(errors.Extend("tracker.RestartAutoPrintGroup()", err))
	}
	tGroup.restartTicker()
	return nil
}

// findGroup takes a trackerGroup name and, if found, returns the object
func (t *SuperTracker) findGroup(group ...string) (*trackerGroup, error) {
	if len(group) == 0 {
		group = append(group, "default")
	} else if len(group) > 1 {
		return nil, errors.New("tracker.findGroup", "AddGauge can only have 0 or 1 value")
	}
	tGroup, ok := t.trackerGroups[group[0]]
	if !ok {
		return nil, errors.New("tracker.findGroup", "Couldn't find group "+group[0])
	}
	return tGroup, nil
}

// findTracker takes a tracker name and, if found, returns the object
func (t *SuperTracker) findTracker(tName string) (tracker, error) {
	tracker := &gauge{}
	for _, tg := range t.trackerGroups {
		if tracker, err := tg.findTracker(tName); err == nil {
			return tracker, nil
		}
	}
	return tracker, errors.New("tracker.findTracker()", fmt.Sprintf("tracker %s was not found", tName))
}

// ETA asdf
func (t *SuperTracker) ETA(tracker string) (string, error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return "", errors.Extend("tracker.ETA()", err)
	}
	endMeasure := trckr.getETA()
	return endMeasure, nil
}
