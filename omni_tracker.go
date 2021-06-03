// package tracker

// import (
// 	"fmt"
// 	"time"

// 	"github.com/morrocker/errors"
// 	"github.com/morrocker/log"
// )

// const (
// 	defLapse = 6000 * time.Millisecond
// 	defGroup = "default"
// 	defVis   = true
// )

// type OmniTracker interface {
// 	// Data modification
// 	AddGauge(string, interface{}, ...string) error
// 	ChangeCurr(string, interface{}) error
// 	Current(string, interface{}) error
// 	ChangeAndReturn(string, ...interface{}) (int64, int64, error)
// 	ChangeTotal(string, interface{}) error
// 	Total(string, interface{}) error
// 	Reset(string) error
// 	// Return values
// 	Values(string) (string, string, error)
// 	RawValues(string) (int64, int64, error)
// 	Print(string)
// 	// Utility functions
// 	StartMeasure(string) (func(), error)
// 	ProgressRate(string) (string, error)
// 	TrueProgressRate(string) (string, error)
// 	PrintFunc(func(), ...string) error
// 	UnitsFunc(string, func(int64) string) error
// 	InitSpdRate(string, int) error
// 	ETA(string) (string, error)
// 	// Periodical functions
// 	StartAutoMeasure(string, time.Duration) error
// 	StopAutoMeasure(string) error
// 	StartAutoPrint(time.Duration, ...string) error
// 	StopAutoPrint(...string) error
// 	RestartGroupAutoPrint(...string) error
// }

// // omniTracker base structure that houses all tracks to monitor task progress
// type omniTracker struct {
// 	trackerGroups map[string]*trackerGroup
// }

// // New creates a new omniTracker. The 'default' group is created with it.
// func New() OmniTracker {
// 	tracker := &omniTracker{
// 		trackerGroups: make(map[string]*trackerGroup),
// 	}
// 	tracker.trackerGroups[defGroup] = newGroup(defGroup)
// 	return tracker
// }

// // AddGauge creates a new gauge on the 'default' group
// func (t *omniTracker) AddGauge(tracker string, total interface{}, group ...string) error {
// 	op := "omni_tracker.AddGauge()"
// 	tGroup, err := t.findGroup(group...)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	if err := tGroup.addGauge(tracker, total); err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	return nil
// }

// // ChangeCurr changes a trackers current value by the set amount
// func (t *omniTracker) ChangeCurr(tracker string, value interface{}) (err error) {
// 	op := "omni_tracker.ChangeCurr()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	val64, err := getInt64(value)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	trckr.changeCurrent(val64)
// 	return
// }

// // SetCurr sets a trackers current value to the one given
// func (t *omniTracker) Current(tracker string, value interface{}) error {
// 	op := "omni_tracker.Current()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	val64, err := getInt64(value)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	trckr.setCurrent(val64)
// 	return nil
// }

// func (t *omniTracker) ChangeAndReturn(tracker string, values ...interface{}) (int64, int64, error) {
// 	op := "omni_tracker.ChangeAndReturn()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return 0, 0, errors.Extend(op, err)
// 	}

// 	curr, err := getInt64(values[0])
// 	if err != nil {
// 		return 0, 0, errors.Extend(op, err)
// 	}
// 	tot, err := getInt64(values[1])
// 	if err != nil {
// 		return 0, 0, errors.Extend(op, err)
// 	}
// 	newCurr, newTot := trckr.changeAndReturn(curr, tot)
// 	return newCurr, newTot, nil
// }

// // Reset sets a trackers current and total value to 0
// func (t *omniTracker) Reset(tracker string) error {
// 	op := "omni_tracker.ResetCurr()"
// 	if err := t.Current(tracker, 0); err != nil {
// 		return errors.New(op, err)
// 	}
// 	if err := t.Total(tracker, 0); err != nil {
// 		return errors.New(op, err)
// 	}
// 	return nil
// }

// // ChangeTotal changes a trackers total value by the given amount
// func (t *omniTracker) ChangeTotal(tracker string, value interface{}) error {
// 	op := "omni_tracker.ChangeTotal()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	val64, err := getInt64(value)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	trckr.changeTotal(val64)
// 	return nil
// }

// // Total sets a trackers value to the one given
// func (t *omniTracker) Total(tracker string, value interface{}) error {
// 	op := "omni_tracker.Total()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	val64, err := getInt64(value)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	trckr.setTotal(val64)
// 	return nil
// }

// // Values returns the current and total values from the given tracker
// func (t *omniTracker) Values(tracker string) (current string, total string, err error) {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		err = errors.Extend("omni_tracker.GetValues()", err)
// 		return
// 	}
// 	return trckr.getValues()
// }

// // RawValues returns the current and total values from the given tracker
// func (t *omniTracker) RawValues(tracker string) (current int64, total int64, err error) {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		err = errors.Extend("omni_tracker.GetValues()", err)
// 		return
// 	}
// 	current, total = trckr.getRawValues()
// 	return
// }

// // PrintFunc executes the print function for the 'default' group
// func (t *omniTracker) PrintFunc(f func(), group ...string) error {
// 	tGroup, err := t.findGroup(group...)
// 	if err != nil {
// 		log.Errorln(errors.Extend("omni_tracker.SetGroupPrintFunc()", err))
// 	}
// 	tGroup.printFunc(f)
// 	return nil
// }

// // PrintGroup executes the print function
// func (t *omniTracker) Print(group string) {
// 	tGroup, err := t.findGroup(group)
// 	if err != nil {
// 		log.Errorln(errors.Extend("omni_tracker.PrintGroup()", err))
// 	}
// 	tGroup.print()
// }

// // Status changes the group's status
// func (t *omniTracker) Status(task string, group ...string) error {
// 	tGroup, err := t.findGroup(group...)
// 	if err != nil {
// 		return errors.Extend("omni_tracker.PrintGroup()", err)
// 	}
// 	tGroup.setStatus(task)
// 	return nil
// }

// // StartMeasure starts a parameter measuring, returning a function that ends the measure
// func (t *omniTracker) StartMeasure(tracker string) (func(), error) {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return nil, errors.Extend("omni_tracker.StartMeasure()", err)
// 	}
// 	endMeasure := trckr.spdMeasureStart()
// 	return endMeasure, nil
// }

// // ProgressRate returns the measures average rate without applying the units modifier function
// func (t *omniTracker) ProgressRate(tracker string) (string, error) {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return "", errors.Extend("omni_tracker.GetProgressRate()", err)
// 	}
// 	return trckr.getRate(), nil
// }

// // TrueProgressRate returns the measure average rate applying the units modifier function
// func (t *omniTracker) TrueProgressRate(tracker string) (string, error) {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		err = errors.Extend("omni_tracker.GetTrueProgressRate()", err)
// 		return "", err
// 	}
// 	return trckr.getRate(), nil
// }

// // UnitsFunc sets the given function to be used to modify the tracked untyped units into a specific unit
// func (t *omniTracker) UnitsFunc(tracker string, f func(int64) string) error {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend("omni_tracker.SetProgressFunction()", err)
// 	}
// 	trckr.setUnitsFunc(f)
// 	return nil
// }

// // InitSpdRate TODO STILL THINKING IF having to initialize this separately is a good idea.
// func (t *omniTracker) InitSpdRate(tracker string, n int) error {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend("omni_tracker.IncreaseCurr()", err)
// 	}
// 	trckr.initSpdRate(n)
// 	return nil
// }

// // StartAutoMeasure starts the automeasuring process. The process will take measures for each
// // time.Duration given for the given tracker. TODO automeasure still uses old way to count time
// func (t *omniTracker) StartAutoMeasure(tracker string, tick time.Duration) error {
// 	op := "omni_tracker.StartAutoMeasure()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}

// 	if err := trckr.startAutoMeasure(tick); err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	return nil

// }

// // StopAutoMeasure stops the automeasuring process for the given tracker
// func (t *omniTracker) StopAutoMeasure(tracker string) error {
// 	op := "omni_tracker.StartAutoMeasure()"
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return errors.Extend(op, err)
// 	}

// 	if err := trckr.stopAutoMeasure(); err != nil {
// 		return errors.Extend(op, err)
// 	}
// 	return nil

// }

// // StartGroupAutoPrint starts the autoprint process for the given group
// func (t *omniTracker) StartAutoPrint(d time.Duration, group ...string) error {
// 	tGroup, err := t.findGroup(group...)
// 	if err != nil {
// 		log.Errorln(errors.Extend("omni_tracker.StartAutoPrintGroup()", err))
// 	}
// 	tGroup.startAutoPrint(d)
// 	return nil
// }

// // StopGroupAutoPrint stops the autoprint process for the given group
// func (t *omniTracker) StopAutoPrint(group ...string) error {
// 	tGroup, err := t.findGroup(group...)
// 	if err != nil {
// 		log.Errorln(errors.Extend("omni_tracker.StopAutoPrintGroup()", err))
// 	}
// 	tGroup.stopAutoPrint()
// 	return nil
// }

// // RestartGroupAutoPrint delays the auto-printing for the given group. Also restart's it
// // if it was stopped
// func (t *omniTracker) RestartGroupAutoPrint(group ...string) error {
// 	tGroup, err := t.findGroup(group...)
// 	if err != nil {
// 		log.Errorln(errors.Extend("omni_tracker.RestartAutoPrintGroup()", err))
// 	}
// 	tGroup.restartTicker()
// 	return nil
// }

// // findGroup takes a trackerGroup name and, if found, returns the object
// func (t *omniTracker) findGroup(group ...string) (*trackerGroup, error) {
// 	if len(group) == 0 {
// 		group = append(group, "default")
// 	} else if len(group) > 1 {
// 		return nil, errors.New("omni_tracker.findGroup", "AddGauge can only have 0 or 1 value")
// 	}
// 	tGroup, ok := t.trackerGroups[group[0]]
// 	if !ok {
// 		return nil, errors.New("omni_tracker.findGroup", "Couldn't find group "+group[0])
// 	}
// 	return tGroup, nil
// }

// // findTracker takes a tracker name and, if found, returns the object
// func (t *omniTracker) findTracker(tName string) (tracker, error) {
// 	tracker := &gauge{}
// 	for _, tg := range t.trackerGroups {
// 		if tracker, err := tg.findTracker(tName); err == nil {
// 			return tracker, nil
// 		}
// 	}
// 	return tracker, errors.New("omni_tracker.findTracker()", fmt.Sprintf("tracker %s was not found", tName))
// }

// // ETA asdf
// func (t *omniTracker) ETA(tracker string) (string, error) {
// 	trckr, err := t.findTracker(tracker)
// 	if err != nil {
// 		return "", errors.Extend("omni_tracker.ETA()", err)
// 	}
// 	endMeasure := trckr.getETA()
// 	return endMeasure, nil
// }
