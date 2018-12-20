package mapping

import (
	"errors"
	"log"
	"math"

	"github.com/aphecetche/pigiron/geo"
)

// ErrInvalidPadCID signals that an invalid pad uid was used / returned
var ErrInvalidPadCID = errors.New("invalid pad uid")

// PadCID is a pad identifier, valid for one cathode only.
type PadCID int

// DualSampaID is a DualSampa identifier.
type DualSampaID int

// DualSampaChannel is a DualSampa channel identifier.
type DualSampaChannelID int

// DEID is a detection element identifier.
type DEID int

// CathodeSegmentation represents the mapping of one cathode
// (either bending or non-bending cathode).
type CathodeSegmentation interface {
	DetElemID() DEID
	NofPads() int
	NofDualSampas() int
	DualSampaID(dualSampaIndex int) (DualSampaID, error)
	IsValid(padcid PadCID) bool
	FindPadByFEE(dualSampaID DualSampaID, dualSampaChannel DualSampaChannelID) (PadCID, error)
	FindPadByPosition(x, y float64) (PadCID, error)
	ForEachPad(padHandler func(padcid PadCID))
	ForEachPadInDualSampa(dualSampaID DualSampaID, padHandler func(padcid PadCID))
	PadDualSampaChannel(padcid PadCID) DualSampaChannelID
	PadDualSampaID(padcid PadCID) DualSampaID
	PadPositionX(padcid PadCID) float64
	PadPositionY(padcid PadCID) float64
	PadSizeX(padcid PadCID) float64
	PadSizeY(padcid PadCID) float64
	GetNeighbours(padcid PadCID) []PadCID
	IsBending() bool
	String(padcid PadCID) string
}

// ForEachDetectionElement loops over all detection elements and
// call the detElemIdHandler function for each of them
func ForEachDetectionElement(detElemIDHandler func(deid DEID)) {
	for _, deid := range []DEID{100, 101, 102, 103,
		200, 201, 202, 203, 300,
		301, 302, 303, 400, 401, 402, 403,
		500, 501, 502, 503, 504, 505, 506, 507, 508,
		509, 510, 511, 512, 513, 514, 515, 516, 517,
		600, 601, 602, 603, 604, 605, 606, 607, 608,
		609, 610, 611, 612, 613, 614, 615, 616, 617,
		700, 701, 702, 703, 704, 705, 706, 707, 708, 709, 710, 711, 712,
		713, 714, 715, 716, 717, 718, 719, 720, 721, 722, 723, 724, 725,
		800, 801, 802, 803, 804, 805, 806, 807, 808, 809, 810, 811, 812,
		813, 814, 815, 816, 817, 818, 819, 820, 821, 822, 823, 824, 825,
		900, 901, 902, 903, 904, 905, 906, 907, 908, 909, 910, 911, 912,
		913, 914, 915, 916, 917, 918, 919, 920, 921, 922, 923, 924, 925,
		1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012,
		1013, 1014, 1015, 1016, 1017, 1018, 1019, 1020, 1021, 1022, 1023, 1024, 1025} {
		detElemIDHandler(deid)
	}
}

// ForOneDetectionElementOfEachSegmentationType loops over one detection element per segmentation type
// and call the detElemIdHandler function for each of them
func ForOneDetectionElementOfEachSegmentationType(detElemIDHandler func(deid DEID)) {
	for _, deid := range []DEID{100, 300, 500, 501, 502, 503, 504, 600, 601, 602,
		700, 701, 702, 703, 704, 705, 706, 902, 903, 904, 905} {
		detElemIDHandler(deid)
	}
}

// PlaneAbbreviation returns a short name for a bending/non-bending plane
func PlaneAbbreviation(isBendingPlane bool) string {
	if isBendingPlane {
		return "B"
	}
	return "NB"
}

// ComputeCathodeBBox return the bounding box of the cathode.
func ComputeBBox(cseg CathodeSegmentation) geo.BBox {
	xmin := math.MaxFloat64
	ymin := xmin
	xmax := -xmin
	ymax := -ymin
	cseg.ForEachPad(func(padcid PadCID) {
		bbox := ComputeCathodePadBBox(cseg, padcid)
		xmin = math.Min(xmin, bbox.Xmin())
		xmax = math.Max(xmax, bbox.Xmax())
		ymin = math.Min(ymin, bbox.Ymin())
		ymax = math.Max(ymax, bbox.Ymax())
	})
	bbox, err := geo.NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		panic(err)
	}
	return bbox
}

// ComputeCathodePadBBox returns the bounding box of one pad of the
// given cathode.
func ComputeCathodePadBBox(cseg CathodeSegmentation, padcid PadCID) geo.BBox {
	x := cseg.PadPositionX(padcid)
	y := cseg.PadPositionY(padcid)
	dx := cseg.PadSizeX(padcid) / 2
	dy := cseg.PadSizeY(padcid) / 2
	bbox, err := geo.NewBBox(x-dx, y-dy, x+dx, y+dy)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}
	return bbox
}

// NewSegmentation creates a segmentation object for the given
// detection element plane (aka cathode).
func NewCathodeSegmentation(deid DEID, isBendingPlane bool) CathodeSegmentation {
	segType, err := detElemID2SegType(deid)
	if err != nil {
		return nil
	}
	builder := getCathodeSegmentationBuilder(segType)
	if builder == nil {
		return nil
	}
	seg := builder.Build(isBendingPlane, deid)
	return seg
}
