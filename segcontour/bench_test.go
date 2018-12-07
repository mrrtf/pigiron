package segcontour_test

import (
	"strconv"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/segcontour"
)

func getSegs() map[mapping.DEID]mapping.CathodeSegmentation {
	segs := make(map[mapping.DEID]mapping.CathodeSegmentation)
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		segs[deid] = mapping.NewCathodeSegmentation(deid, true)
	})
	return segs
}

func BenchmarkSegmentationComputeBBoxViaPadLoop(b *testing.B) {

	segs := getSegs()

	for deid, seg := range segs {
		b.Run(strconv.Itoa(int(deid)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mapping.ComputeBBox(seg)
			}
		})
	}
}

func BenchmarkSegmentationComputeBBoxViaContour(b *testing.B) {

	segs := getSegs()

	for deid, seg := range segs {
		b.Run(strconv.Itoa(int(deid)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = segcontour.BBox(seg)
			}
		})
	}
}
