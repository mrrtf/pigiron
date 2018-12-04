package segcontour_test

import (
	"strconv"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/segcontour"
)

func getSegs() map[int]mapping.CathodeSegmentation {
	segs := make(map[int]mapping.CathodeSegmentation)
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		segs[detElemID] = mapping.NewCathodeSegmentation(detElemID, true)
	})
	return segs
}

func BenchmarkSegmentationComputeBBoxViaPadLoop(b *testing.B) {

	segs := getSegs()

	for detElemID, seg := range segs {
		b.Run(strconv.Itoa(detElemID), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mapping.ComputeBBox(seg)
			}
		})
	}
}

func BenchmarkSegmentationComputeBBoxViaContour(b *testing.B) {

	segs := getSegs()

	for detElemID, seg := range segs {
		b.Run(strconv.Itoa(detElemID), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = segcontour.BBox(seg)
			}
		})
	}
}
