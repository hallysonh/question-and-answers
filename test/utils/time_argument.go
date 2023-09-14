package test_utils

import (
	"database/sql/driver"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AlmostNow struct {
	AcceptableDiffInSec float64
}

func NewAlmostNow(acceptableDiffInSec float64) *AlmostNow {
	return &AlmostNow{AcceptableDiffInSec: acceptableDiffInSec}
}

// Match satisfies sqlmock.Argument interface
func (a AlmostNow) Match(v driver.Value) bool {
	timeValue, ok := v.(time.Time)
	if !ok {
		return false
	}
	diffInSec := time.Now().Sub(timeValue).Abs().Seconds()
	return diffInSec < a.AcceptableDiffInSec
}
