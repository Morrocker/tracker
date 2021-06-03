package tracker

import "time"

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
