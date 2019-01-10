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

// BenchmarkNeighboursArray checks the cost of getting the neighbours
// of some pads, and also checks that the overhead of using the segmentation
// (instead of cathode segmentations) is minimal.
func BenchmarkNeighboursArray(b *testing.B) {
	var deid mapping.DEID = 100
	catsegB := mapping.NewCathodeSegmentation(deid, true)
	catsegNB := mapping.NewCathodeSegmentation(deid, false)
	seg := mapping.NewSegmentation(deid)
	nei := make([]int, 13)
	b.Run("Cathode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			catsegB.ForEachPad(func(padcid mapping.PadCID) {
				_ = catsegB.GetNeighboursArray(padcid, nei)
			})
			catsegNB.ForEachPad(func(padcid mapping.PadCID) {
				_ = catsegNB.GetNeighboursArray(padcid, nei)
			})
		}
	})
	b.Run("Segmentation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			seg.ForEachPad(func(paduid mapping.PadUID) {
				_ = seg.GetNeighboursArray(paduid, nei)
			})
		}
	})
}
