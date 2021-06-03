package tracker

// import (
// 	"fmt"
// 	"time"

// 	"github.com/morrocker/errors"
// 	"github.com/morrocker/log"
// )

// type trackerGroup struct {
// 	trackers   map[string]tracker
// 	ticksLapse time.Duration
// 	status     string
// 	ticker     *time.Ticker
// 	prntFunc   func()
// }

// // NewGroup creates a new tracker group within a SuperTracker.
// func newGroup(name string) *trackerGroup {
// 	group := &trackerGroup{
// 		trackers: make(map[string]tracker),
// 	}
// 	return group
// }

// // AddGauge call AddGaugeOn using the default group
// func (g *trackerGroup) addGauge(trackerName string, total interface{}) error {
// 	op := "tracker_group.addGauge()"
// 	if _, err := g.findTracker(trackerName); err == nil {
// 		return errors.New(op, fmt.Sprintf("tracker name %s already taken", trackerName))
// 	}
// 	total64, err := getInt64(total)
// 	if err != nil {
// 		return errors.New(op, err)
// 	}
// 	g.trackers[trackerName] = newGauge(trackerName, total64)
// 	return nil
// }

// func (g *trackerGroup) findTracker(t string) (tracker, error) {
// 	tracker := &gauge{}
// 	for name, tracker := range g.trackers {
// 		if t == name {
// 			return tracker, nil
// 		}
// 	}
// 	err := errors.New("tracker.findTracker()", "Did not find tracker "+t)
// 	return tracker, err
// }

// func (g *trackerGroup) printFunc(f func()) {
// 	g.prntFunc = f
// }

// func (g *trackerGroup) print() {
// 	if g.prntFunc == nil {
// 		log.Errorln(errors.New("tracker_group.print()", "print function is not set!"))
// 		return
// 	}
// 	g.prntFunc()
// }

// func (g *trackerGroup) setStatus(s string) {
// 	g.status = s
// }

// func (g *trackerGroup) startAutoPrint(d time.Duration) {
// 	if err := g.checkTicker(); err == nil {
// 		g.restartTicker()
// 	}
// 	g.ticksLapse = d
// 	g.ticker = time.NewTicker(d)

// 	go func() {
// 		for range g.ticker.C {
// 			g.print()
// 		}
// 	}()
// }

// func (g *trackerGroup) stopAutoPrint() error {
// 	if err := g.checkTicker(); err != nil {
// 		return errors.Extend("tracker_group.stopAutoPrint()", err)
// 	}
// 	g.ticker.Stop()
// 	return nil
// }

// func (g *trackerGroup) restartTicker() error {
// 	if g.ticker == nil {
// 		return errors.New("tracker_group.resetTicker()", "Ticker hasn't been set")
// 	}
// 	g.ticker.Reset(g.ticksLapse)
// 	return nil
// }

// func (g *trackerGroup) checkTicker() error {
// 	if g.ticker == nil {
// 		return errors.New("tracker_group.checkTicker()", "Ticker hasn't been set")
// 	}
// 	return nil
// }
