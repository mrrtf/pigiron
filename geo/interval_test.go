package geo

import "testing"

func TestIntervalCtorMustErrorIfBeginIsAfterEnd(t *testing.T) {
	_, err := newInterval(24, 3)
	if err == nil {
		t.Error("interval ctor should error if begin > end")
	}
}

func TestIntervalCtorThrowsIfBeginEqualsEnd(t *testing.T) {
	_, err := newInterval(24.24, 24.24)
	if err == nil {
		t.Error("interval ctor should error if begin == end")
	}
}

func TestIsFullyContained(t *testing.T) {
	i, err := newInterval(0.01, 0.05)
	if err != nil {
		t.Error("interval could not be created while it should have been")
	}
	a := interval{0, 0.04}
	if a.isFullyContainedIn(i) != false {
		t.Error("was expecting false here")
	}
	b := interval{0.01, 0.02}
	if b.isFullyContainedIn(i) != true {
		t.Error("was expecting true here")
	}
}

func TestExtend(t *testing.T) {
	i, _ := newInterval(0.01, 0.05)
	j, _ := newInterval(0.05, 0.07)

	expected, _ := newInterval(0.01, 0.07)

	ok := i.extend(j)
	if !ok || !equalInterval(i, expected) {
		t.Error("extend failed")
	}

	k, _ := newInterval(0.05, 0.10)
	ok = i.extend(k)
	if ok {
		t.Error("this extend should have failed")
	}
}

func TestIntervalString(t *testing.T) {
	expected := "[0.01,0.05]"
	i, _ := newInterval(0.01, 0.05)
	s := i.String()
	if s != expected {
		t.Errorf("expected string:%s and got:%s", expected, s)
	}
}
