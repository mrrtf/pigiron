package mapping_test

import (
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

func TestNewPadGroupType(t *testing.T) {

	pgt := mapping.NewPadGroupType(
		2, 48, []int{15, 16, 14, 17, 13, 18, 12, 19, 11, 20, 10, 21, 9, 22, 8, 23, 7, 24, 6, 25, 5, 26, 4, 27,
			3, 28, 2, 29, 1, 30, 0, 31, -1, 48, -1, 49, -1, 50, -1, 51, -1, 52, -1, 53, -1, 54, -1, 55,
			-1, 56, -1, 57, -1, 58, -1, 59, -1, 60, -1, 61, -1, 62, -1, 63, -1, 32, -1, 33, -1, 34, -1, 35,
			-1, 36, -1, 37, -1, 38, -1, 39, -1, 40, -1, 41, -1, 42, -1, 43, -1, 44, -1, 45, -1, 46, -1, 47},
	)

	if pgt.NofPads() != 64 {
		t.Error("padgrouptype should have 64 pads")
	}

	pgt = mapping.NewPadGroupType(
		28, 2, []int{47, 45, 43, 41, 39, 37, 35, 33, 63, 61, 59, 57, 55, 53, 51, 49, 31, 29, 27,
			25, 23, 21, 19, 17, 15, 13, 11, 9, 46, 44, 42, 40, 38, 36, 34, 32, 62, 60,
			58, 56, 54, 52, 50, 48, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8})

	if pgt.NofPads() != 56 {
		t.Error("padgrouptype should have 64 pads")
	}
}
