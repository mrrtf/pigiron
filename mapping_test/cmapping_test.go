package mapping_test

import (
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

var (
	seg = mapping.NewSegmentation(100, true)
)

func TestNewCSegmentationMustNotErrorIfDetElemIdIsValid(t *testing.T) {

	seg := mapping.NewSegmentation(100, true)
	if seg == nil {
		t.Fatalf("Could not create segmentation")
	}
}

func TestNewCSegmentationMustErrorIfDetElemIdIsNotValid(t *testing.T) {

	seg := mapping.NewSegmentation(-1, true)
	if seg != nil {
		t.Fatalf("Should have failed here")
	}
	seg = mapping.NewSegmentation(121, true)
	if seg != nil {
		t.Fatalf("Should have failed here")
	}
}
func TestValidFindPadByFEE(t *testing.T) {
	_, err := seg.FindPadByFEE(102, 3)
	if err != nil {
		t.Errorf("Should get a valid pad here")
	}
}
func TestInvalidFindPadByFEE(t *testing.T) {
	_, err := seg.FindPadByFEE(214, 14)
	if err == nil {
		t.Errorf("Should not get a valid pad here")
	}
}
func TestMustErrorIfDualSampaChannelIsNotBetween0And63(t *testing.T) {
	_, err := seg.FindPadByFEE(102, -1)
	if err == nil {
		t.Errorf("Should _not_ get a valid pad here")
	}
	_, err = seg.FindPadByFEE(102, 64)
	if err == nil {
		t.Errorf("Should _not_ get a valid pad here")
	}
}

func TestPositionOfOnePadInDE100Bending(t *testing.T) {
	p1, err := seg.FindPadByFEE(76, 9)
	if err != nil {
		t.Error("Should get a valid pad")
	}
	p2, err := seg.FindPadByPosition(1.575, 18.69)
	if err != nil {
		t.Error("Should get a valid pad")
	}
	if p1 != p2 {
		t.Error("Should get the same pads here")
	}
}
