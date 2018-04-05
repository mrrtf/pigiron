package mapping_test

import (
	"log"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

var (
	seg = mapping.NewSegmentation(100, true)
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

/*

BOOST_AUTO_TEST_CASE(DualSampasWithLessThan64Pads)
{
  std::map<int, int> non64;
  forOneDetectionElementOfEachSegmentationType([&non64](int detElemId) {
    for (auto plane : { true, false }) {
      Segmentation seg{ detElemId, plane };
      for (int i = 0; i < seg.nofDualSampas(); ++i) {
        int n{ 0 };
        seg.forEachPadInDualSampa(seg.dualSampaId(i), [&n](int paduid) { ++n; });
        if (n != 64) {
          non64[n]++;
        }
      }
    }
  });

  BOOST_CHECK_EQUAL(non64[31], 1);
  BOOST_CHECK_EQUAL(non64[32], 2);
  BOOST_CHECK_EQUAL(non64[39], 1);
  BOOST_CHECK_EQUAL(non64[40], 3);
  BOOST_CHECK_EQUAL(non64[46], 2);
  BOOST_CHECK_EQUAL(non64[48], 10);
  BOOST_CHECK_EQUAL(non64[49], 1);
  BOOST_CHECK_EQUAL(non64[50], 1);
  BOOST_CHECK_EQUAL(non64[52], 3);
  BOOST_CHECK_EQUAL(non64[54], 2);
  BOOST_CHECK_EQUAL(non64[55], 3);
  BOOST_CHECK_EQUAL(non64[56], 114);
  BOOST_CHECK_EQUAL(non64[57], 3);
  BOOST_CHECK_EQUAL(non64[58], 2);
  BOOST_CHECK_EQUAL(non64[59], 1);
  BOOST_CHECK_EQUAL(non64[60], 6);
  BOOST_CHECK_EQUAL(non64[62], 4);
  BOOST_CHECK_EQUAL(non64[63], 7);

  int n{ 0 };
  for (auto p : non64) {
    n += p.second;
  }

  BOOST_CHECK_EQUAL(n, 166);
}
*/

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
