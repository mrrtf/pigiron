package mapping

import (
	"fmt"
	"log"
	"math"

	"github.com/mrrtf/pigiron/geo"
)

type PadUID int

// Segmentation is the main entry point to the MCH mapping.
// It encompasses both cathodes of a detection element.
//
type Segmentation interface {
	DetElemID() DEID
	NofPads() int
	NofDualSampas() int
	IsValid(paduid PadUID) bool
	IsBendingPad(paduid PadUID) bool
	FindPadByFEE(dualSampaID DualSampaID, dualSampaChannel DualSampaChannelID) (PadUID, error)
	FindPadPairByPosition(x, y float64) (PadUID, PadUID, error)
	ForEachPad(padHandler func(paduid PadUID))
	ForEachPadInDualSampa(dualSampaID DualSampaID, padHandler func(paduid PadUID))
	PadDualSampaChannel(paduid PadUID) DualSampaChannelID
	PadDualSampaID(paduid PadUID) DualSampaID
	PadPositionX(paduid PadUID) float64
	PadPositionY(paduid PadUID) float64
	PadSizeX(paduid PadUID) float64
	PadSizeY(paduid PadUID) float64
	GetNeighbourIDs(paduid PadUID, neighbours []int) int
	Bending() CathodeSegmentation
	NonBending() CathodeSegmentation
	String(paduid PadUID) string
}

type segmentation struct {
	bending    CathodeSegmentation
	nonBending CathodeSegmentation
	// PadUID handled by Segmentation ranges from 0..NofPads-1
	// Convention is that the first padids are for bending
	// cathode (up to nof pads of that cathode) and then
	// for non-bending
	padUIDOffset int
}

var InvalidPadUID PadUID = -1

// NewSegmentation creates a Segmentation object for the given
// detection element.
// The returned segmentation spans both cathodes (bending and non-bending).
func NewSegmentation(deid DEID) Segmentation {
	bseg := NewCathodeSegmentation(deid, true)
	if bseg == nil {
		return nil
	}
	nbseg := NewCathodeSegmentation(deid, false)
	if nbseg == nil {
		return nil
	}
	seg := &segmentation{bending: bseg, nonBending: nbseg}
	seg.padUIDOffset = seg.bending.NofPads()
	return seg
}

func (seg *segmentation) Bending() CathodeSegmentation {
	return seg.bending
}

func (seg *segmentation) NonBending() CathodeSegmentation {
	return seg.nonBending
}

func (seg *segmentation) IsBendingPad(paduid PadUID) bool {
	return int(paduid) < seg.padUIDOffset
}

func (seg *segmentation) getCathSeg(paduid PadUID) (CathodeSegmentation, PadCID, error) {
	if paduid < 0 {
		return nil, 0, fmt.Errorf("invalid pad uid")
	}
	if int(paduid) < seg.padUIDOffset {
		return seg.bending, PadCID(paduid), nil
	}
	return seg.nonBending, PadCID(int(paduid) - seg.padUIDOffset), nil
}

func (seg *segmentation) String(paduid PadUID) string {
	catseg, padcid, err := seg.getCathSeg(paduid)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("PAD %10d ", paduid) + catseg.String(padcid)
}

func (seg *segmentation) DetElemID() DEID {
	return seg.bending.DetElemID()
}

func (seg *segmentation) NofPads() int {
	return seg.bending.NofPads() + seg.nonBending.NofPads()
}

func (seg *segmentation) NofDualSampas() int {
	return seg.bending.NofDualSampas() + seg.nonBending.NofDualSampas()
}

func (seg *segmentation) IsValid(paduid PadUID) bool {
	cseg, p, err := seg.getCathSeg(paduid)
	if err != nil {
		return false
	}
	return cseg.IsValid(p)
}

func (seg *segmentation) padC2UID(padcid PadCID, isBending bool) PadUID {
	if isBending {
		return PadUID(padcid)
	}
	return PadUID(int(padcid) + seg.padUIDOffset)
}

func (seg *segmentation) FindPadByFEE(dualSampaID DualSampaID, dualSampaChannel DualSampaChannelID) (PadUID, error) {
	var padcid PadCID
	var err error
	var isBending bool = dualSampaID < 1024
	if isBending {
		padcid, err = seg.bending.FindPadByFEE(dualSampaID, dualSampaChannel)
	} else {
		padcid, err = seg.nonBending.FindPadByFEE(dualSampaID, dualSampaChannel)
	}
	if err != nil {
		return InvalidPadUID, err
	}
	return seg.padC2UID(padcid, isBending), nil
}

func (seg *segmentation) FindPadPairByPosition(x, y float64) (PadUID, PadUID, error) {
	bp, erb := seg.bending.FindPadByPosition(x, y)
	nbp, ernb := seg.nonBending.FindPadByPosition(x, y)
	var err error
	if erb != nil {
		err = erb
	}
	if ernb != nil {
		if err == nil {
			err = erb
		} else {
			err = fmt.Errorf("%s and %s", erb.Error(), ernb.Error())
		}
	}
	return seg.padC2UID(bp, true), seg.padC2UID(nbp, false), err
}

func f2cuid(padHandler func(paduid PadUID), offset int) func(padcid PadCID) {
	return func(padcid PadCID) {
		padHandler(PadUID(padcid) + PadUID(offset))
	}
}

func (seg *segmentation) ForEachPad(padHandler func(paduid PadUID)) {
	seg.bending.ForEachPad(f2cuid(padHandler, 0))
	seg.nonBending.ForEachPad(f2cuid(padHandler, seg.padUIDOffset))
}

func (seg *segmentation) ForEachPadInDualSampa(dualSampaID DualSampaID, padHandler func(paduid PadUID)) {
	if dualSampaID < 1024 {
		seg.bending.ForEachPad(f2cuid(padHandler, 0))
	} else {
		seg.bending.ForEachPad(f2cuid(padHandler, seg.padUIDOffset))
	}
}

func (seg *segmentation) GetNeighbourIDs(paduid PadUID, neighbours []int) int {
	cseg, p, err := seg.getCathSeg(paduid)
	if err != nil {
		log.Fatalf("Should not happen")
	}
	n := cseg.GetNeighbourIDs(p, neighbours)
	if !cseg.IsBending() {
		for i := 0; i < n; i++ {
			neighbours[i] += seg.padUIDOffset
		}
	}
	return n
}

func (seg *segmentation) PadDualSampaChannel(paduid PadUID) DualSampaChannelID {
	cseg, p, _ := seg.getCathSeg(paduid)
	return cseg.PadDualSampaChannel(p)
}

func (seg *segmentation) PadDualSampaID(paduid PadUID) DualSampaID {
	cseg, p, _ := seg.getCathSeg(paduid)
	return cseg.PadDualSampaID(p)
}

func (seg *segmentation) PadPositionX(paduid PadUID) float64 {
	cseg, p, _ := seg.getCathSeg(paduid)
	return cseg.PadPositionX(p)
}

func (seg *segmentation) PadPositionY(paduid PadUID) float64 {
	cseg, p, _ := seg.getCathSeg(paduid)
	return cseg.PadPositionY(p)
}

func (seg *segmentation) PadSizeX(paduid PadUID) float64 {
	cseg, p, _ := seg.getCathSeg(paduid)
	return cseg.PadSizeX(p)
}

func (seg *segmentation) PadSizeY(paduid PadUID) float64 {
	cseg, p, _ := seg.getCathSeg(paduid)
	return cseg.PadSizeY(p)
}

// ComputeSegmentationBBox return the bounding box of the
// detection element represented by it segmentation.
func ComputeSegmentationBBox(seg Segmentation) geo.BBox {
	xmin := math.MaxFloat64
	ymin := xmin
	xmax := -xmin
	ymax := -ymin
	var xpadmin, ypadmin, xpadmax, ypadmax float64
	seg.ForEachPad(func(paduid PadUID) {
		ComputePadBBox(seg, paduid, &xpadmin, &ypadmin, &xpadmax, &ypadmax)
		xmin = math.Min(xmin, xpadmin)
		xmax = math.Max(xmax, xpadmax)
		ymin = math.Min(ymin, ypadmin)
		ymax = math.Max(ymax, ypadmax)
	})
	bbox, err := geo.NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		panic(err)
	}
	return bbox
}

// ComputePadBBox fills the coordinates (xmin,ymin,xmax,ymax) of the bounding
// box of one pad of the given segmentation.
func ComputePadBBox(padps PadSizerPositioner, paduid PadUID, xmin, ymin, xmax, ymax *float64) {
	x := padps.PadPositionX(paduid)
	y := padps.PadPositionY(paduid)
	dx := padps.PadSizeX(paduid) / 2
	dy := padps.PadSizeY(paduid) / 2
	*xmin = x - dx
	*xmax = x + dx
	*ymin = y - dy
	*ymax = y + dy
}
