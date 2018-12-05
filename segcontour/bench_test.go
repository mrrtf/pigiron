package segcontour_test

import (
	"strconv"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/segcontour"
)

func getSegs() map[int]mapping.CathodeSegmentation {
	segs := make(map[int]mapping.CathodeSegmentation)
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid int) {
		segs[deid] = mapping.NewCathodeSegmentation(deid, true)
	})
	return segs
}

func BenchmarkSegmentationComputeBBoxViaPadLoop(b *testing.B) {

	segs := getSegs()

	for deid, seg := range segs {
		b.Run(strconv.Itoa(deid), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mapping.ComputeBBox(seg)
			}
		})
	}
}

func BenchmarkSegmentationComputeBBoxViaContour(b *testing.B) {

	segs := getSegs()

	for deid, seg := range segs {
		b.Run(strconv.Itoa(deid), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = segcontour.BBox(seg)
			}
		})
	}
}
