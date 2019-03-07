package mapping_test

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

func BenchmarkSegmentationCreationPerDE(b *testing.B) {

	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		b.Run(strconv.Itoa(int(deid)), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = mapping.NewCathodeSegmentation(deid, true)
				_ = mapping.NewCathodeSegmentation(deid, false)
			}
		})
	})
}

func BenchmarkSegmentationCreation(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
			_ = mapping.NewCathodeSegmentation(deid, true)
			_ = mapping.NewCathodeSegmentation(deid, false)
		})
	}
}

type SegPair map[bool]mapping.CathodeSegmentation

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
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		const n = 100000
		for _, isBendingPlane := range []bool{true, false} {
			b.Run(fmt.Sprintf("findPadByPositions(%d,%v)", deid, isBendingPlane), func(b *testing.B) {
				seg := mapping.NewCathodeSegmentation(deid, isBendingPlane)
				bbox := mapping.ComputeBBox(seg)
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

type DC struct {
	D mapping.DualSampaID
	C mapping.DualSampaChannelID
}

var (
	detElemIds []mapping.DEID
)

func init() {

	detElemIds = []mapping.DEID{100, 300, 501, 1025}
}

func BenchmarkByFEE(b *testing.B) {
	for _, deid := range detElemIds {
		for _, isBendingPlane := range []bool{true, false} {
			planeName := "B"
			if isBendingPlane == false {
				planeName = "NB"
			}
			seg := mapping.NewCathodeSegmentation(deid, isBendingPlane)
			var dcs []DC
			seg.ForEachPad(func(padcid mapping.PadCID) {
				dcs = append(dcs, DC{D: seg.PadDualSampaID(padcid), C: seg.PadDualSampaChannel(padcid)})
			})
			b.Run(strconv.Itoa(int(deid))+planeName, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					for _, pad := range dcs {
						seg.FindPadByFEE(pad.D, pad.C)
					}
				}
			})
		}
	}
}

// BenchmarkNeighbourIDs checks the cost of getting the neighbours
// of some pads, and also checks that the overhead of using the segmentation
// (instead of cathode segmentations) is minimal.
func BenchmarkNeighbourIDs(b *testing.B) {
	var deid mapping.DEID = 100
	catsegB := mapping.NewCathodeSegmentation(deid, true)
	catsegNB := mapping.NewCathodeSegmentation(deid, false)
	seg := mapping.NewSegmentation(deid)
	nei := make([]int, 13)
	b.Run("Cathode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			catsegB.ForEachPad(func(padcid mapping.PadCID) {
				_ = catsegB.GetNeighbourIDs(padcid, nei)
			})
			catsegNB.ForEachPad(func(padcid mapping.PadCID) {
				_ = catsegNB.GetNeighbourIDs(padcid, nei)
			})
		}
	})
	b.Run("Segmentation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			seg.ForEachPad(func(paduid mapping.PadUID) {
				_ = seg.GetNeighbourIDs(paduid, nei)
			})
		}
	})
}

var iresult int

// BenchmarkForEachPad compares the time it takes to look over
// all pads using the two functions ForEachPad and ForEachPadInArea,
// where the latter is given an infinite box as a parameter.
// Note that in this case the ForEachPadInArea is expected to be
// (much) slower than ForEachPad, as it's a limiting case (i.e.
// if you want to loop over all pads, use ForEachPad, but if you want
// to loop over a -hopefully rather small- area, use ForEachPadInArea.
// But at least we should get the same number of pads...
func BenchmarkForEachPad(b *testing.B) {
	for _, deid := range []mapping.DEID{100, 500} {
		seg := mapping.NewSegmentation(deid)
		b.Run(fmt.Sprintf("ForEachPadDE%d", deid), func(b *testing.B) {
			var n int
			for i := 0; i < b.N; i++ {
				n = 0
				seg.ForEachPad(func(paduid mapping.PadUID) {
					n++
				})
			}
			iresult = n
		})
		b.Run(fmt.Sprintf("ForEachPadInAreaDE%d", deid), func(b *testing.B) {
			var n int
			xmin := -math.MaxFloat64
			ymin := -math.MaxFloat64
			xmax := math.MaxFloat64
			ymax := math.MaxFloat64
			for i := 0; i < b.N; i++ {
				n = 0
				seg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(paduid mapping.PadUID) {
					n++
				})
			}
			iresult = n
		})
	}

}
