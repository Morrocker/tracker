package tracker

import (
	"fmt"
	"time"

	"github.com/morrocker/errors"
	"github.com/morrocker/log"
)

const (
	defSeparator = "|"
	defLapse     = 6000 * time.Millisecond
	defLineMode  = "singleline"
	defGaugeMode = "division"
	defGroup     = "default"
	defVis       = true
)

// SuperTracker base structure that houses all tracks to monitor task progress
type SuperTracker struct {
	trackerGroups map[string]*trackerGroup
}

type tracker interface {
	setMode(string)
	setCurrent(int64)
	setTotal(int64)
	changeCurrent(int64)
	changeTotal(int64)
	getRawValues() (int64, int64)
	getValues() (string, string, error)
	print() string
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
func (t *SuperTracker) AddGauge(tracker, printName string, total interface{}) error {
	t.AddGaugeOn(defGroup, tracker, printName, total)
	return nil
}

// AddGaugeOn creates and adds a gauge type tracker to a tracker group.
func (t *SuperTracker) AddGaugeOn(group, tracker, printName string, total interface{}) error {
	op := "tracker.AddGaugeOn()"
	tGroup, err := t.findGroup(group)
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

// IncreaseCurr increases a trackers current value by 1
func (t *SuperTracker) IncreaseCurr(tracker string) error {
	return t.ChangeCurr(tracker, 1)
}

// SetCurr sets a trackers current value to the one given
func (t *SuperTracker) SetCurr(tracker string, value interface{}) error {
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

// Reset sets a trackers current and total value to 0
func (t *SuperTracker) Reset(tracker string) error {
	op := "tracker.ResetCurr()"
	if err := t.SetCurr(tracker, 0); err != nil {
		return errors.New(op, err)
	}
	if err := t.Total(tracker, 0); err != nil {
		return errors.New(op, err)
	}
	return nil
}

// ResetCurr sets a trackers current value to 0
func (t *SuperTracker) ResetCurr(tracker string) (err error) {
	return t.SetCurr(tracker, 0)
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

// Mode change the tracker mode
func (t *SuperTracker) Mode(tracker, mode string) (err error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend("tracker.SetMode()", err)
	}
	trckr.setMode(mode)
	return
}

// PrintFunc executes the print function for the 'default' group
func (t *SuperTracker) PrintFunc(f func()) error {
	return t.GroupPrintFunc(defGroup, f)
}

// GroupPrintFunc executes the print function for the given group
func (t *SuperTracker) GroupPrintFunc(group string, f func()) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.SetGroupPrintFunc()", err))
	}
	tGroup.printFunc(f)
	return nil
}

// Print prints progress of all visible trackers
func (t *SuperTracker) Print() {
	t.PrintGroup(defGroup)
}

// PrintGroup prints progress of all visible trackers
func (t *SuperTracker) PrintGroup(group string) {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.PrintGroup()", err))
	}
	tGroup.print()
}

// Status changes the 'default' group's status
func (t *SuperTracker) Status(group, task string) error {
	return t.GroupStatus(defGroup, task)
}

// GroupStatus changes the given group's status
func (t *SuperTracker) GroupStatus(group, task string) error {
	tGroup, err := t.findGroup(group)
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

// StartAutoPrint starts the autoprint process for the 'default' group
func (t *SuperTracker) StartAutoPrint(d time.Duration) (err error) {
	err = t.StartGroupAutoPrint(defGroup, d)
	return
}

// StartGroupAutoPrint starts the autoprint process for the given group
func (t *SuperTracker) StartGroupAutoPrint(group string, d time.Duration) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.StartAutoPrintGroup()", err))
	}
	tGroup.startAutoPrint(d)
	return nil
}

// StopAutoPrint stops the autoprint process for the 'default' group
func (t *SuperTracker) StopAutoPrint() (err error) {
	err = t.StopGroupAutoPrint(defGroup)
	return
}

// StopGroupAutoPrint stops the autoprint process for the given group
func (t *SuperTracker) StopGroupAutoPrint(group string) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.StopAutoPrintGroup()", err))
	}
	tGroup.stopAutoPrint()
	return nil
}

// RestartAutoPrint delays the auto-printing for the 'default' group. Also restart's it
// if it was stopped
func (t *SuperTracker) RestartAutoPrint() (err error) {
	err = t.RestartGroupAutoPrint(defGroup)
	return
}

// RestartGroupAutoPrint delays the auto-printing for the given group. Also restart's it
// if it was stopped
func (t *SuperTracker) RestartGroupAutoPrint(group string) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.RestartAutoPrintGroup()", err))
	}
	tGroup.restartTicker()
	return nil
}

// findGroup takes a trackerGroup name and, if found, returns the object
func (t *SuperTracker) findGroup(name string) (*trackerGroup, error) {
	group, ok := t.trackerGroups[name]
	if !ok {
		return nil, errors.New("tracker.findGroup", "Couldn't find group "+name)
	}
	return group, nil
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
