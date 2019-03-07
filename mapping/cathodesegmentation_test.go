package mapping_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aphecetche/pigiron/mapping"

	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func TestCathodeNumberOfDetectionElementIs156(t *testing.T) {
	nde := 0
	mapping.ForEachDetectionElement(func(deid mapping.DEID) {
		nde++
	})

	if nde != 156 {
		t.Errorf("Expected 156 detection elements, got %d", nde)
	}
}

func TestCathodeNewSegmentationMustNotErrorIfDetElemIdIsValid(t *testing.T) {

	cseg := mapping.NewCathodeSegmentation(100, true)
	if cseg == nil {
		t.Fatalf("Could not create cathodesegmentation")
	}
}

func TestCathodeNewSegmentationMustErrorIfDetElemIdIsNotValid(t *testing.T) {

	cseg := mapping.NewCathodeSegmentation(-1, true)
	if cseg != nil {
		t.Fatalf("Should have failed here")
	}
	cseg = mapping.NewCathodeSegmentation(121, true)
	if cseg != nil {
		t.Fatalf("Should have failed here")
	}
}

var testcathodedeid = []mapping.DEID{100, 300, 500, 501, 502, 503, 504, 600, 601, 602, 700,
	701, 702, 703, 704, 705, 706, 902, 903, 904, 905}

func TestCathodeCreateSegmentation(t *testing.T) {
	for _, de := range testcathodedeid {
		for _, plane := range []bool{true, false} {
			cseg := mapping.NewCathodeSegmentation(de, plane)
			if cseg == nil {
				t.Fatalf("could not create cathode segmentation for DE %d", de)
			}
		}
	}
}

func TestCathodeNofPads(t *testing.T) {
	var tv = []struct {
		de                mapping.DEID
		nofBendingPads    int
		nofNonBendingPads int
	}{
		{100, 14392, 14280},
		{300, 13947, 13986},
		{902, 4480, 3136},
		{702, 4160, 2912},
		{701, 4096, 2880},
		{601, 3648, 2560},
		{501, 3568, 2496},
		{602, 3200, 2240},
		{700, 3200, 2240},
		{502, 3120, 2176},
		{600, 3008, 2112},
		{500, 2928, 2048},
		{903, 2880, 2016},
		{703, 2560, 1792},
		{904, 2240, 1568},
		{503, 1920, 1344},
		{704, 1920, 1344},
		{504, 1280, 896},
		{905, 1280, 896},
		{705, 960, 672},
		{706, 640, 448},
	}

	for _, tt := range tv {
		b := mapping.NewCathodeSegmentation(tt.de, true)
		nb := mapping.NewCathodeSegmentation(tt.de, false)
		if b.NofPads() != tt.nofBendingPads {
			t.Errorf("DE %d : expected %d pads in bending plane. Got %d", tt.de, b.NofPads(), tt.nofBendingPads)
		}
		if nb.NofPads() != tt.nofNonBendingPads {
			t.Errorf("DE %d : Expected %d pads in non bending plane. Got %d", tt.de, nb.NofPads(), tt.nofNonBendingPads)
		}
	}
}

func TestCathodeTotalNofBendingFECInSegTypes(t *testing.T) {
	var tv = []struct {
		plane         bool
		nofDualSampas int
	}{
		{false, 1019},
		{true, 1246},
	}

	for _, tt := range tv {
		n := 0
		for _, de := range testcathodedeid {
			cseg := mapping.NewCathodeSegmentation(de, tt.plane)
			n += cseg.NofDualSampas()
		}
		if n != tt.nofDualSampas {
			t.Errorf("Expected %d got %d", tt.nofDualSampas, n)
		}
	}
}

func TestCathodeNofFEC(t *testing.T) {
	var tv = []struct {
		de                      mapping.DEID
		nofBendingDualSampas    int
		nofNonBendingDualSampas int
	}{
		{100, 226, 225},
		{300, 221, 222},
		{902, 70, 50},
		{702, 65, 46},
		{701, 64, 46},
		{601, 57, 40},
		{501, 56, 39},
		{602, 50, 35},
		{700, 50, 36},
		{502, 49, 34},
		{600, 47, 33},
		{500, 46, 32},
		{903, 45, 33},
		{703, 40, 29},
		{904, 35, 26},
		{503, 30, 21},
		{704, 30, 22},
		{504, 20, 14},
		{905, 20, 16},
		{705, 15, 12},
		{706, 10, 8},
	}

	for _, tt := range tv {
		b := mapping.NewCathodeSegmentation(tt.de, true)
		nb := mapping.NewCathodeSegmentation(tt.de, false)
		if b.NofDualSampas() != tt.nofBendingDualSampas {
			t.Errorf("DE %d : expected %d dual sampas in bending plane. Got %d", tt.de, b.NofDualSampas(), tt.nofBendingDualSampas)
		}
		if nb.NofDualSampas() != tt.nofNonBendingDualSampas {
			t.Errorf("DE %d : Expected %d dual sampas in non bending plane. Got %d", tt.de, nb.NofDualSampas(), tt.nofNonBendingDualSampas)
		}
	}
}

func TestCathodeNofPadsInSegmentations(t *testing.T) {
	npads := 0
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		for _, plane := range []bool{true, false} {
			cseg := mapping.NewCathodeSegmentation(deid, plane)
			if cseg == nil {
				log.Fatalf("Got nil seg for detElemId %d plane %v", deid, plane)
			}
			npads += cseg.NofPads()
		}

	})
	if npads != 143469 {
		t.Errorf("Expected 143469 pads : got %d", npads)
	}
}

func TestCathodeNofSegmentations(t *testing.T) {
	n := 0
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		n += 2
	})
	if n != 42 {
		t.Errorf("Expected 42 segmentations : got %d", n)
	}
}

func TestCathodeDualSampasWithLessThan64Pads(t *testing.T) {

	non64 := make(map[int]int)
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		for _, plane := range []bool{true, false} {
			cseg := mapping.NewCathodeSegmentation(deid, plane)
			for i := 0; i < cseg.NofDualSampas(); i++ {
				n := 0
				dualSampaID, _ := cseg.DualSampaID(i)
				cseg.ForEachPadInDualSampa(dualSampaID, func(paduid mapping.PadCID) {
					n++
				})
				if n != 64 {
					non64[n]++
				}
			}
		}
	})

	var expected = []struct {
		npads     int
		occurence int
	}{
		{31, 1},
		{32, 2},
		{39, 1},
		{40, 3},
		{46, 2},
		{48, 10},
		{49, 1},
		{50, 1},
		{52, 3},
		{54, 2},
		{55, 3},
		{56, 114},
		{57, 3},
		{58, 2},
		{59, 1},
		{60, 6},
		{62, 4},
		{63, 7},
	}

	for _, tt := range expected {
		if non64[tt.npads] != tt.occurence {
			t.Errorf("Expected %d dual sampas with %d pads, but got %d", tt.occurence, tt.npads, non64[tt.npads])
		}
	}

	n := 0
	for _, v := range non64 {
		n += v
	}

	if n != 166 {
		t.Errorf("Expected 166 dual sampas with a number of pads different from 64 and got %d", n)
	}
}

func TestCathodeMustErrorIfDualSampaChannelIsNotBetween0And63(t *testing.T) {
	cseg := mapping.NewCathodeSegmentation(100, true)
	_, err := cseg.FindPadByFEE(102, -1)
	if err == nil {
		t.Errorf("Should _not_ get a valid pad here")
	}
	_, err = cseg.FindPadByFEE(102, 64)
	if err == nil {
		t.Errorf("Should _not_ get a valid pad here")
	}
}

func TestCathodePositionOfOnePadInDE100Bending(t *testing.T) {
	cseg := mapping.NewCathodeSegmentation(100, true)
	p1, err := cseg.FindPadByFEE(76, 9)
	if err != nil {
		t.Errorf("Should get a valid pad: %v", err)
	}
	p2, err := cseg.FindPadByPosition(1.575, 18.69)
	if err != nil {
		t.Errorf("Should get a valid pad: %v", err)
	}
	if p1 != p2 {
		t.Errorf("Should get the same pads here p1=%v p2=%v", p1, p2)
		PrintCathodePad(os.Stdout, cseg, p1)
		PrintCathodePad(os.Stdout, cseg, p2)
	}
}

func TestCathodeValidFindPadByFEE(t *testing.T) {
	cseg := mapping.NewCathodeSegmentation(100, true)
	_, err := cseg.FindPadByFEE(102, 3)
	if err != nil {
		t.Errorf("Should get a valid pad here")
	}
}
func TestCathodeInvalidFindPadByFEE(t *testing.T) {
	cseg := mapping.NewCathodeSegmentation(100, true)
	_, err := cseg.FindPadByFEE(214, 14)
	if err == nil {
		t.Errorf("Should not get a valid pad here")
	}
}

type Point struct {
	x, y float64
}

//TODO: implement me ?
func checkGaps(t *testing.T, cseg *mapping.CathodeSegmentation) []Point {
	return []Point{}
}

func dumpToFile(filename string, cseg *mapping.CathodeSegmentation, points []Point) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString("toto")
	w.Flush()
}

func TestCathodeNoGapWithinPads(t *testing.T) {
	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		for _, isBendingPlane := range []bool{true, false} {
			cseg := mapping.NewCathodeSegmentation(deid, isBendingPlane)
			g := checkGaps(t, &cseg)
			if len(g) != 0 {
				dumpToFile(fmt.Sprintf("bug-gap-%v-%s.html",
					deid,
					mapping.PlaneAbbreviation(isBendingPlane)),
					&cseg, g)
			}
		}
	})
}

// FIXME: uncomment those tests
//
// func TestCathodeForEachPadCountPads(t *testing.T) {
// 	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
// 		for _, b := range []bool{true, false} {
// 			cseg := mapping.NewCathodeSegmentation(deid, b)
// 			npads := 0
// 			cseg.ForEachPad(func(paduid mapping.PadCID) {
// 				npads++
// 			})
// 			if npads != cseg.NofPads() {
// 				t.Errorf("DE %v isBending %v : expected %v pads but got %v from ForEachPad loop", deid, b, cseg.NofPads(), npads)
// 			}
// 		}
// 	})
// }
//
// func TestCathodeForEachPadInAreaCountPads(t *testing.T) {
// 	// mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
// 	for _, deid := range []mapping.DEID{100, 706} {
// 		for _, b := range []bool{false} {
// 			cseg := mapping.NewCathodeSegmentation(deid, b)
// 			npads := 0
// 			xmin := -math.MaxFloat64
// 			xmax := math.MaxFloat64
// 			ymin := -math.MaxFloat64
// 			ymax := math.MaxFloat64
// 			cseg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(paduid mapping.PadCID) {
// 				npads++
// 			})
// 			if npads != cseg.NofPads() {
// 				t.Errorf("DE %v isBending %v : expected %v pads but got %v from ForEachPad loop", deid, b, cseg.NofPads(), npads)
// 			}
// 		}
// 	}
// 	// })
// }
//
// func TestCathodeForEachPadInAreaBug(t *testing.T) {
// 	cseg := mapping.NewCathodeSegmentation(100, true)
// 	x := 52.5
// 	y := 20.0
// 	dx := 5.0
// 	dy := 1.0
//
// 	xmin := x - dx/2.0
// 	xmax := x + dx/2.0
// 	ymin := y - dy/2.0
// 	ymax := y + dy/2.0
//
// 	padcids := make(map[int]struct{})
//
// 	cseg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(padcid mapping.PadCID) {
// 		padcids[int(padcid)] = struct{}{}
// 	})
//
// 	expected := []int{
// 		2975,
// 		2991,
// 		4496,
// 		3038,
// 		2974,
// 		3070,
// 		3134,
// 		3023,
// 		3071,
// 		3055,
// 		2990,
// 		4480,
// 		4528,
// 		4592,
// 		4512,
// 		3022,
// 		3039,
// 		2958,
// 		2959,
// 		3135,
// 		3054,
// 		1,
// 		2,
// 	}
// 	if len(padcids) != len(expected) {
// 		t.Errorf("Want %d pads - Got %d", len(expected), len(padcids))
// 		for p := range padcids {
// 			fmt.Println(p, cseg.PadDualSampaID(mapping.PadCID(p)), cseg.PadDualSampaChannel(mapping.PadCID(p)))
// 		}
// 	}
// }
//
// func BenchmarkCathodeForEachPadInAreaBug(b *testing.B) {
// 	cseg := mapping.NewCathodeSegmentation(100, true)
// 	x := 52.5
// 	y := 20.0
// 	dx := 5.0
// 	dy := 1.0
//
// 	xmin := x - dx/2.0
// 	xmax := x + dx/2.0
// 	ymin := y - dy/2.0
// 	ymax := y + dy/2.0
//
// 	padcids := make(map[int]struct{})
//
// 	for i := 0; i < b.N; i++ {
// 		cseg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(padcid mapping.PadCID) {
// 			padcids[int(padcid)] = struct{}{}
// 		})
// 	}
// }
//
// func TestCathodeForEachPadInAreaDev(t *testing.T) {
// 	cseg := mapping.NewCathodeSegmentation(100, true)
// 	x := 52.5
// 	y := 20.0
// 	dx := 5.0
// 	dy := 1.0
//
// 	xmin := x - dx/2.0
// 	xmax := x + dx/2.0
// 	ymin := y - dy/2.0
// 	ymax := y + dy/2.0
//
// 	padcids := make(map[int]struct{})
//
// 	err := forEachPadInArea(cseg, xmin, ymin, xmax, ymax, func(padcid mapping.PadCID) {
// 		padcids[int(padcid)] = struct{}{}
// 	})
//
// 	if err != nil {
// 		t.Errorf("forEachPadInArea failed")
// 	}
// 	expected := []int{
// 		2975,
// 		2991,
// 		4496,
// 		3038,
// 		2974,
// 		3070,
// 		3134,
// 		3023,
// 		3071,
// 		3055,
// 		2990,
// 		4480,
// 		4528,
// 		4592,
// 		4512,
// 		3022,
// 		3039,
// 		2958,
// 		2959,
// 		3135,
// 		3054,
// 		1,
// 		2,
// 	}
// 	if len(padcids) != len(expected) {
// 		t.Errorf("Want %d pads - Got %d", len(expected), len(padcids))
// 		for p := range padcids {
// 			fmt.Println(p, cseg.PadDualSampaID(mapping.PadCID(p)), cseg.PadDualSampaChannel(mapping.PadCID(p)))
// 		}
// 	}
// }
//
// // test method to develop foreachpadinarea algo
// func forEachPadInArea(cseg mapping.CathodeSegmentation, xmin, ymin, xmax, ymax float64, padfunc func(padcid mapping.PadCID)) error {
//
// 	bbox, err := geo.NewBBox(xmin, ymin, xmax, ymax)
// 	if err != nil {
// 		log.Fatal("oups")
// 	}
// 	csegContour := segcontour.Contour(cseg)
// 	overlap, err := geo.ClipPolygon(csegContour[0], bbox)
// 	if err != nil {
// 		return err
// 	}
// 	obox := overlap.BBox()
// 	fmt.Printf("%v\nBBOX=%v\n", overlap, obox)
// 	micron := 1E-4 // distances are in cm in segmentation. This converts to micrometers
// 	delta := 1.0 * micron
// 	x := obox.Xmin() + delta
// 	y := obox.Ymin() + delta
//
// 	// startingPoint, err := cseg.FindPadByPosition(x, y)
// 	// if err != nil {
// 	// 	log.Fatalf("Could not get a pad at pos %v,%v !", x, y)
// 	// }
// 	// fmt.Printf("%6d ", int(startingPoint))
// 	// PrintCathodePad(os.Stdout, cseg, startingPoint)
//
// 	ypadsize := 0.0
//
// 	for y < obox.Ymax() {
// 		xmin, xmax := geo.GetPolygonXRangeAtY(overlap, y)
// 		x = xmin + delta
// 		for x < xmax {
// 			pad, err := cseg.FindPadByPosition(x, y)
// 			if err != nil {
// 				break
// 			}
// 			fmt.Printf("%6d ", int(pad))
// 			PrintCathodePad(os.Stdout, cseg, pad)
// 			x += cseg.PadPositionX(pad) + cseg.PadSizeX(pad) + delta
// 			if cseg.PadSizeY(pad) != ypadsize {
// 				fmt.Printf("should handle ypadsize change from %v to %v\n", ypadsize, cseg.PadSizeY(pad))
// 				ypadsize = cseg.PadSizeY(pad)
// 			}
// 		}
// 		y += ypadsize
// 	}
// 	return nil
// }
