package mapping_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

func BenchmarkSegmentationCreationPerDE(b *testing.B) {

	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		b.Run(strconv.Itoa(detElemID), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = mapping.NewSegmentation(detElemID, true)
				_ = mapping.NewSegmentation(detElemID, false)
			}
		})
	})
}

func BenchmarkSegmentationCreation(b *testing.B) {

	for i := 0; i < b.N; i++ {
		mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
			_ = mapping.NewSegmentation(detElemID, true)
			_ = mapping.NewSegmentation(detElemID, false)
		})
	}
}

type SegPair map[bool]mapping.Segmentation

type TestPoint struct {
	x, y float64
}

func generateUniformTestPoints(n int, box geo.BBox) []TestPoint {
	var testpoints = make([]TestPoint, n)
	for i := 0; i < n; i++ {
		x := box.Xmin() + rand.Float64()*box.Width()
		y := box.Ymin() + rand.Float64()*box.Height()
		testpoints[i] = TestPoint{x, y}
	}
	return testpoints
}

func BenchmarkPositions(b *testing.B) {
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		const n = 100000
		for _, isBendingPlane := range []bool{true, false} {
			b.Run(fmt.Sprintf("findPadByPositions(%d,%v)", detElemID, isBendingPlane), func(b *testing.B) {
				seg := mapping.NewSegmentation(detElemID, isBendingPlane)
				bbox := mapping.ComputeBbox(seg)
				testpoints := generateUniformTestPoints(n, bbox)
				for i := 0; i < b.N; i++ {
					for _, tp := range testpoints {
						seg.FindPadByPosition(tp.x, tp.y)
					}
				}
			})
		}
	})
}

	// seg.ForEachPad(func(paduid int) {
	// 	x := seg.PadPositionX(paduid)
	// 	y := seg.PadPositionY(paduid)
	// 	dx := seg.PadSizeX(paduid) / 2
	// 	dy := seg.PadSizeY(paduid) / 2
	// 	xmin = math.Min(xmin, x-dx)
	// 	xmax = math.Max(xmax, x+dx)
	// 	ymin = math.Min(ymin, y-dy)
	// 	ymax = math.Max(ymax, y+dy)
	// })

	type DC struct {
		D int
		C int
	}
func BenchmarkByFEE(b *testing.B) {

    dc := []DC{}

    seg := mapping.NewSegmentation(100,true)

	seg.ForEachPad(func(paduid int){
		dc = append(dc,DC{D:seg.PadDualSampaID(paduid),C:seg.PadDualSampaChannel(paduid)})
	})

	for i := 0; i < b.N; i++ {
		for _,pad := range dc {
			seg.FindPadByFEE(pad.D,pad.C)
	}
	}

}

