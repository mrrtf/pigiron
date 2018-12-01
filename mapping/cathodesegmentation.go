package mapping

import (
	"fmt"
	"io"
	"math"

	"github.com/aphecetche/pigiron/geo"
)

// Segmentation is the main entry point to the MCH mapping
type CathodeSegmentation interface {
	DetElemID() int
	NofPads() int
	NofDualSampas() int
	DualSampaID(dualSampaIndex int) (int, error)
	IsValid(padid int) bool
	FindPadByFEE(dualSampaID, dualSampaChannel int) (int, error)
	FindPadByPosition(x, y float64) (int, error)
	ForEachPad(padHandler func(paduid int))
	ForEachPadInDualSampa(dualSampaID int, padHandler func(paduid int))
	PadDualSampaChannel(paduid int) int
	PadDualSampaID(paduid int) int
	PadPositionX(paduid int) float64
	PadPositionY(paduid int) float64
	PadSizeX(paduid int) float64
	PadSizeY(paduid int) float64
	GetNeighbours(paduid int) []int
	IsBending() bool
	setDetElemID(de int)
}

// ForEachDetectionElement loops over all detection elements and call the detElemIdHandler function
// for each of them
func ForEachDetectionElement(detElemIDHandler func(detElemID int)) {
	for _, detElemID := range []int{100, 101, 102, 103, 200, 201, 202, 203, 300, 301, 302, 303, 400, 401, 402, 403, 500, 501,
		502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 512, 513, 514, 515, 516, 517, 600, 601,
		602, 603, 604, 605, 606, 607, 608, 609, 610, 611, 612, 613, 614, 615, 616, 617, 700, 701,
		702, 703, 704, 705, 706, 707, 708, 709, 710, 711, 712, 713, 714, 715, 716, 717, 718, 719,
		720, 721, 722, 723, 724, 725, 800, 801, 802, 803, 804, 805, 806, 807, 808, 809, 810, 811,
		812, 813, 814, 815, 816, 817, 818, 819, 820, 821, 822, 823, 824, 825, 900, 901, 902, 903,
		904, 905, 906, 907, 908, 909, 910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 920, 921,
		922, 923, 924, 925, 1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013,
		1014, 1015, 1016, 1017, 1018, 1019, 1020, 1021, 1022, 1023, 1024, 1025} {
		detElemIDHandler(detElemID)
	}
}

// ForOneDetectionElementOfEachSegmentationType loops over one detection element per segmentation type
// and call the detElemIdHandler function for each of them
func ForOneDetectionElementOfEachSegmentationType(detElemIDHandler func(detElemID int)) {
	for _, detElemID := range []int{100, 300, 500, 501, 502, 503, 504, 600, 601, 602, 700, 701, 702, 703, 704, 705, 706, 902, 903, 904, 905} {
		detElemIDHandler(detElemID)
	}
}

// PlaneAbbreviation returns a short name for a bending/non-bending plane
func PlaneAbbreviation(isBendingPlane bool) string {
	if isBendingPlane {
		return "B"
	}
	return "NB"
}

// PrintPad prints all known information about a pad
func PrintPad(out io.Writer, cseg CathodeSegmentation, paduid int) {
	if !cseg.IsValid(paduid) {
		fmt.Printf("invalid pad")
		return
	}
	fmt.Fprintf(out, "DE %4d DSID %4d CH %2d X %7.2f Y %7.2f DX %7.2f DY %7.2f\n",
		cseg.DetElemID(),
		cseg.PadDualSampaID(paduid),
		cseg.PadDualSampaChannel(paduid),
		cseg.PadPositionX(paduid),
		cseg.PadPositionY(paduid),
		cseg.PadSizeX(paduid),
		cseg.PadSizeY(paduid))

}

func ComputeBbox(cseg CathodeSegmentation) geo.BBox {
	xmin := math.MaxFloat64
	ymin := xmin
	xmax := -xmin
	ymax := -ymin
	cseg.ForEachPad(func(paduid int) {
		x := cseg.PadPositionX(paduid)
		y := cseg.PadPositionY(paduid)
		dx := cseg.PadSizeX(paduid) / 2
		dy := cseg.PadSizeY(paduid) / 2
		xmin = math.Min(xmin, x-dx)
		xmax = math.Max(xmax, x+dx)
		ymin = math.Min(ymin, y-dy)
		ymax = math.Max(ymax, y+dy)
	})
	bbox, err := geo.NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		panic(err)
	}
	return bbox
}
