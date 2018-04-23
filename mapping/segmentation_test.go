package mapping_test

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

var (
// seg = mapping.NewSegmentation(100, true)
)

func TestNumberOfDetectionElementIs156(t *testing.T) {
	nde := 0
	mapping.ForEachDetectionElement(func(detElemID int) {
		nde++
	})

	if nde != 156 {
		t.Errorf("Expected 156 detection elements, got %d", nde)
	}
}

func TestNewSegmentationMustNotErrorIfDetElemIdIsValid(t *testing.T) {

	seg := mapping.NewSegmentation(100, true)
	if seg == nil {
		t.Fatalf("Could not create segmentation")
	}
}

func TestNewSegmentationMustErrorIfDetElemIdIsNotValid(t *testing.T) {

	seg := mapping.NewSegmentation(-1, true)
	if seg != nil {
		t.Fatalf("Should have failed here")
	}
	seg = mapping.NewSegmentation(121, true)
	if seg != nil {
		t.Fatalf("Should have failed here")
	}
}

var detElemIDs = []int{100, 300, 500, 501, 502, 503, 504, 600, 601, 602, 700,
	701, 702, 703, 704, 705, 706, 902, 903, 904, 905}

func TestCreateSegmentation(t *testing.T) {
	for _, de := range detElemIDs {
		for _, plane := range []bool{true, false} {
			seg := mapping.NewSegmentation(de, plane)
			if seg == nil {
				t.Fatalf("could not create segmentation for DE %d", de)
			}
		}
	}
}

func TestNofPads(t *testing.T) {
	var tv = []struct {
		de                int
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
		b := mapping.NewSegmentation(tt.de, true)
		nb := mapping.NewSegmentation(tt.de, false)
		if b.NofPads() != tt.nofBendingPads {
			t.Errorf("DE %d : expected %d pads in bending plane. Got %d", tt.de, b.NofPads(), tt.nofBendingPads)
		}
		if nb.NofPads() != tt.nofNonBendingPads {
			t.Errorf("DE %d : Expected %d pads in non bending plane. Got %d", tt.de, nb.NofPads(), tt.nofNonBendingPads)
		}
	}
}

func TestTotalNofBendingFECInSegTypes(t *testing.T) {
	var tv = []struct {
		plane         bool
		nofDualSampas int
	}{
		{false, 1019},
		{true, 1246},
	}

	for _, tt := range tv {
		n := 0
		for _, de := range detElemIDs {
			seg := mapping.NewSegmentation(de, tt.plane)
			n += seg.NofDualSampas()
		}
		if n != tt.nofDualSampas {
			t.Errorf("Expected %d got %d", tt.nofDualSampas, n)
		}
	}
}

func TestNofFEC(t *testing.T) {
	var tv = []struct {
		de                      int
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
		b := mapping.NewSegmentation(tt.de, true)
		nb := mapping.NewSegmentation(tt.de, false)
		if b.NofDualSampas() != tt.nofBendingDualSampas {
			t.Errorf("DE %d : expected %d dual sampas in bending plane. Got %d", tt.de, b.NofDualSampas(), tt.nofBendingDualSampas)
		}
		if nb.NofDualSampas() != tt.nofNonBendingDualSampas {
			t.Errorf("DE %d : Expected %d dual sampas in non bending plane. Got %d", tt.de, nb.NofDualSampas(), tt.nofNonBendingDualSampas)
		}
	}
}

func TestNofPadsInSegmentations(t *testing.T) {
	npads := 0
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		for _, plane := range []bool{true, false} {
			seg := mapping.NewSegmentation(detElemID, plane)
			if seg == nil {
				log.Fatalf("Got nil seg for detElemId %d plane %v", detElemID, plane)
			}
			npads += seg.NofPads()
		}

	})
	if npads != 143469 {
		t.Errorf("Expected 143469 pads : got %d", npads)
	}
}

func TestNofSegmentations(t *testing.T) {
	n := 0
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		n += 2
	})
	if n != 42 {
		t.Errorf("Expected 42 segmentations : got %d", n)
	}
}

func TestDualSampasWithLessThan64Pads(t *testing.T) {

	non64 := make(map[int]int)
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		for _, plane := range []bool{true, false} {
			seg := mapping.NewSegmentation(detElemID, plane)
			for i := 0; i < seg.NofDualSampas(); i++ {
				n := 0
				dualSampaID, _ := seg.DualSampaID(i)
				seg.ForEachPadInDualSampa(dualSampaID, func(paduid int) {
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

func TestMustErrorIfDualSampaChannelIsNotBetween0And63(t *testing.T) {
	seg := mapping.NewSegmentation(100, true)
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
	seg := mapping.NewSegmentation(100, true)
	p1, err := seg.FindPadByFEE(76, 9)
	if err != nil {
		t.Errorf("Should get a valid pad: %v", err)
	}
	p2, err := seg.FindPadByPosition(1.575, 18.69)
	if err != nil {
		t.Errorf("Should get a valid pad: %v", err)
	}
	if p1 != p2 {
		t.Errorf("Should get the same pads here p1=%v p2=%v", p1, p2)
		mapping.PrintPad(os.Stdout, seg, p1)
		mapping.PrintPad(os.Stdout, seg, p2)
	}
}

func TestValidFindPadByFEE(t *testing.T) {
	seg := mapping.NewSegmentation(100, true)
	_, err := seg.FindPadByFEE(102, 3)
	if err != nil {
		t.Errorf("Should get a valid pad here")
	}
}
func TestInvalidFindPadByFEE(t *testing.T) {
	seg := mapping.NewSegmentation(100, true)
	_, err := seg.FindPadByFEE(214, 14)
	if err == nil {
		t.Errorf("Should not get a valid pad here")
	}
}

type Point struct {
	x, y float64
}

func checkGaps(t *testing.T, seg *mapping.Segmentation) []Point {
	return []Point{}
}

func dumpToFile(filename string, seg *mapping.Segmentation, points []Point) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString("toto")
	w.Flush()
}

func TestNoGapWithinPads(t *testing.T) {
	mapping.ForOneDetectionElementOfEachSegmentationType(func(detElemID int) {
		for _, isBendingPlane := range []bool{true, false} {
			seg := mapping.NewSegmentation(detElemID, isBendingPlane)
			g := checkGaps(t, &seg)
			if len(g) != 0 {
				dumpToFile(fmt.Sprintf("bug-gap-%v-%s.html",
					detElemID,
					mapping.PlaneAbbreviation(isBendingPlane)),
					&seg, g)
			}
		}
	})
}
