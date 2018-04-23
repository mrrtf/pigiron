package segcontour_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aphecetche/pigiron/segcontour"

	"github.com/aphecetche/pigiron/mapping"
)

func getSegs() map[int]mapping.Segmentation {
	fmt.Println("getSegs")
	segs := make(map[int]mapping.Segmentation)
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		segs[detElemID] = mapping.NewSegmentation(detElemID, true)
	})
	return segs
}

func BenchmarkSegmentationComputeBBoxViaPadLoop(b *testing.B) {

	segs := getSegs()

	for detElemID, seg := range segs {
		b.Run(strconv.Itoa(detElemID), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mapping.ComputeBbox(seg)
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
