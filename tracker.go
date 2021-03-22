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
	setOrder(int)
	setMode(string)
	setCurrent(int64)
	setTotal(int64)
	changeCurrent(int64)
	changeTotal(int64)
	getOrder() int
	getRawValues() (int64, int64)
	getValues() (string, string, error)
	print() string
	spdMeasureStart() func()
	setUnitsFunc(func(int64) string)
	getRawRate() int64
	getRate() string
	getETA() string
	initSpdRate(int)
	startAutoMeasure(int) error
	stopAutoMeasure() error
}

// New creates a new SuperTracker. The default group is created with it.
func New() *SuperTracker {
	tracker := &SuperTracker{
		trackerGroups: make(map[string]*trackerGroup),
	}
	tracker.trackerGroups[defGroup] = newGroup(defGroup)
	return tracker
}

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

// SetLineMode calls SetGroupLineMode using the default group
func (t *SuperTracker) SetLineMode(mode string) (err error) {
	return t.SetGroupLineMode(defGroup, mode)
}

// SetGroupLineMode changes a groups line mode to the set value
func (t *SuperTracker) SetGroupLineMode(group, mode string) error {
	op := "tracker.SetGroupLineMode()"
	tGroup, err := t.findGroup(group)
	if err != nil {
		return errors.Extend(op, err)
	}
	tGroup.setLineMode("mode")
	return nil
}

func (t *SuperTracker) SetEtaTracker(tracker string) error {
	if err := t.SetGroupEtaTracker(defGroup, tracker); err != nil {
		return err
	}
	return nil
}

func (t *SuperTracker) SetGroupEtaTracker(group, tracker string) error {
	op := "tracker.SetGroupEtaTracker"
	tGroup, err := t.findGroup(group)
	if err != nil {
		return errors.Extend(op, err)
	}
	if err := tGroup.setEtaTracker(tracker); err != nil {
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

// ResetCurr sets a trackers current value to 0
func (t *SuperTracker) Reset(tracker string) error {
	op := "tracker.ResetCurr()"
	if err := t.SetCurr(tracker, 0); err != nil {
		return errors.New(op, err)
	}
	if err := t.SetTotal(tracker, 0); err != nil {
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

// SetTotal sets a trackers value to the one given
func (t *SuperTracker) SetTotal(tracker string, value interface{}) error {
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

// GetValues asdfasfd
func (t *SuperTracker) GetValues(tracker string) (current string, total string, err error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		err = errors.Extend("tracker.GetValues()", err)
		return
	}
	return trckr.getValues()
}

// GetValues asdfasfd
func (t *SuperTracker) GetRawValues(tracker string) (current int64, total int64, err error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		err = errors.Extend("tracker.GetValues()", err)
		return
	}
	current, total = trckr.getRawValues()
	return
}

// SetMode asfasf
func (t *SuperTracker) SetMode(tracker, mode string) (err error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend("tracker.SetMode()", err)
	}
	trckr.setMode(mode)
	return
}

func (t *SuperTracker) SetPrintFunc(group string, f func() string) error {
	return t.SetGroupPrintFunc(defGroup, f)
}

func (t *SuperTracker) SetGroupPrintFunc(group string, f func() string) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.SetGroupPrintFunc()", err))
	}
	tGroup.setPrintFunc(f)
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

// func (t *SuperTracker) setPrintOrder(group string) {
// 	unordered := make(map[string]int)
// 	ordered := []string{}

// 	for key, tracker := range t.trackerGroups[group].trackers {
// 		unordered[key] = tracker.getOrder()
// 	}
// 	for {
// 		min := 0
// 		minTracker := ""
// 		for trackerName, orderVal := range unordered {
// 			if min == 0 {
// 				min = orderVal
// 				minTracker = trackerName
// 			} else if min > orderVal {
// 				min = orderVal
// 				minTracker = trackerName
// 			}
// 		}
// 		ordered = append(ordered, minTracker)
// 		delete(unordered, minTracker)
// 		if len(unordered) == 0 {
// 			break
// 		}
// 	}
// 	t.trackerGroups[group].order = ordered
// }

// SetTask sdfs asfa // REWORK THIS
func (t *SuperTracker) SetTask(group, task string) {
	t.trackerGroups[group].format.status = task
}

func (t *SuperTracker) StartMeasure(tracker string) (func(), error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return nil, errors.Extend("tracker.StartMeasure()", err)
	}
	endMeasure := trckr.spdMeasureStart()
	return endMeasure, nil
}

func (t *SuperTracker) GetProgressRate(tracker string) (string, error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return "", errors.Extend("tracker.GetProgressRate()", err)
	}
	return trckr.getRate(), nil
}

func (t *SuperTracker) GetTrueProgressRate(tracker string) (string, error) {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		err = errors.Extend("tracker.GetTrueProgressRate()", err)
		return "", err
	}
	return trckr.getRate(), nil
}

func (t *SuperTracker) SetUnitsFunc(tracker string, f func(int64) string) error {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend("tracker.SetProgressFunction()", err)
	}
	trckr.setUnitsFunc(f)
	return nil
}

func (t *SuperTracker) InitSpdRate(tracker string, n int) error {
	trckr, err := t.findTracker(tracker)
	if err != nil {
		return errors.Extend("tracker.IncreaseCurr()", err)
	}
	trckr.initSpdRate(n)
	return nil
}

func (t *SuperTracker) StartAutoMeasure(tracker string, tick int) error {
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

func (t *SuperTracker) StartAutoPrint(d time.Duration) (err error) {
	err = t.StartAutoPrintGroup(defGroup, d)
	return
}

func (t *SuperTracker) StartAutoPrintGroup(group string, d time.Duration) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.StartAutoPrintGroup()", err))
	}
	tGroup.startAutoPrint(d)
	return nil
}

func (t *SuperTracker) StopAutoPrint() (err error) {
	err = t.StopAutoPrintGroup(defGroup)
	return
}

func (t *SuperTracker) StopAutoPrintGroup(group string) error {
	tGroup, err := t.findGroup(group)
	if err != nil {
		log.Errorln(errors.Extend("tracker.StopAutoPrintGroup()", err))
	}
	tGroup.stopAutoPrint()
	return nil
}
func (t *SuperTracker) RestartAutoPrint() (err error) {
	err = t.RestartAutoPrintGroup(defGroup)
	return
}

func (t *SuperTracker) RestartAutoPrintGroup(group string) error {
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
