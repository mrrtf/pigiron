package mapping_test

import (
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

var (
	seg = mapping.NewSegmentation(100, true)
)

type deIncrementer struct {
	nde int
}

func (de *deIncrementer) Do(detElemID int) {
	de.nde++
}

func TestNumberOfDetectionElementIs156(t *testing.T) {
	var deinc deIncrementer
	mapping.ForEachDetectionElement(&deinc)
	if deinc.nde != 156 {
		t.Errorf("Expected 156 detection elements, got %d", deinc.nde)
	}
}

func TestNewCSegmentationMustNotErrorIfDetElemIdIsValid(t *testing.T) {

	seg := mapping.NewSegmentation(100, true)
	if seg == nil {
		t.Fatalf("Could not create segmentation")
	}
}

func TestNewCSegmentationMustErrorIfDetElemIdIsNotValid(t *testing.T) {

	seg := mapping.NewSegmentation(-1, true)
	if seg != nil {
		t.Fatalf("Should have failed here")
	}
	seg = mapping.NewSegmentation(121, true)
	if seg != nil {
		t.Fatalf("Should have failed here")
	}
}
func TestValidFindPadByFEE(t *testing.T) {
	_, err := seg.FindPadByFEE(102, 3)
	if err != nil {
		t.Errorf("Should get a valid pad here")
	}
}
func TestInvalidFindPadByFEE(t *testing.T) {
	_, err := seg.FindPadByFEE(214, 14)
	if err == nil {
		t.Errorf("Should not get a valid pad here")
	}
}
func TestMustErrorIfDualSampaChannelIsNotBetween0And63(t *testing.T) {
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
	p1, err := seg.FindPadByFEE(76, 9)
	if err != nil {
		t.Error("Should get a valid pad")
	}
	p2, err := seg.FindPadByPosition(1.575, 18.69)
	if err != nil {
		t.Error("Should get a valid pad")
	}
	if p1 != p2 {
		t.Error("Should get the same pads here")
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
