package geo

import (
	"errors"
	"fmt"
)

// in order to reduce the amount of code (and because the usage of this struct
// is limited to this package), we are not going through the hassle of defining
// an interface just to avoid the usage of bare interval{a,b} (where a could be
// bigger than b). We just assume we're good citizen within this package
// and just use newInterval function whenever we need a (properly constructed) interval
type interval struct {
	b, e float64
}

func newInterval(b, e float64) (interval, error) {
	if b > e || EqualFloat(b, e) {
		return interval{}, errors.New("begin should be strictly < end")
	}
	return interval{b, e}, nil
}

func (i *interval) begin() float64 {
	return i.b
}

func (i *interval) end() float64 {
	return i.e
}

func (i *interval) isFullyContainedIn(j interval) bool {
	return (j.begin() < i.begin() || EqualFloat(j.begin(), i.begin())) &&
		(i.end() < j.end() || EqualFloat(j.end(), i.end()))
}

func (i *interval) extend(j interval) bool {
	if EqualFloat(j.begin(), i.end()) {
		i.e = j.end()
		return true
	}
	return false
}

func (i *interval) String() string {
	return fmt.Sprintf("[%v,%v]", i.begin(), i.end())
}

func equalInterval(i, j interval) bool {
	return EqualFloat(i.begin(), j.begin()) &&
		EqualFloat(i.end(), j.end())
}
