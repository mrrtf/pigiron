package mapping

import (
	"fmt"
	"log"
	"math"

	"github.com/aphecetche/pigiron/geo"
)

type PadUID int

// Segmentation is the main entry point to the MCH mapping.
// It encompasses both cathodes of a detection element.
//
type Segmentation interface {
	DetElemID() int
	NofPads() int
	NofDualSampas() int
	IsValid(paduid PadUID) bool
	IsBendingPad(paduid PadUID) bool
	FindPadByFEE(dualSampaID DualSampaID, dualSampaChannel int) (PadUID, error)
	FindPadPairByPosition(x, y float64) (PadUID, PadUID, error)
	ForEachPad(padHandler func(paduid PadUID))
	ForEachPadInDualSampa(dualSampaID DualSampaID, padHandler func(paduid PadUID))
	PadDualSampaChannel(paduid PadUID) int
	PadDualSampaID(paduid PadUID) DualSampaID
	PadPositionX(paduid PadUID) float64
	PadPositionY(paduid PadUID) float64
	PadSizeX(paduid PadUID) float64
	PadSizeY(paduid PadUID) float64
	GetNeighbours(paduid PadUID) []PadUID
	Bending() CathodeSegmentation
	NonBending() CathodeSegmentation
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
func NewSegmentation(deid int) Segmentation {
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
	return seg.nonBending, PadCID(int(paduid) - seg.padUIDOffset - 1), nil
}

func (seg *segmentation) DetElemID() int {
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

func (seg *segmentation) FindPadByFEE(dualSampaID DualSampaID, dualSampaChannel int) (PadUID, error) {
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

func f2cuid(padHandler func(paduid PadUID)) func(padcid PadCID) {
	return func(padcid PadCID) {
		padHandler(PadUID(padcid))
	}
}

func (seg *segmentation) ForEachPad(padHandler func(paduid PadUID)) {
	seg.bending.ForEachPad(f2cuid(padHandler))
	seg.nonBending.ForEachPad(f2cuid(padHandler))
}

func (seg *segmentation) ForEachPadInDualSampa(dualSampaID DualSampaID, padHandler func(paduid PadUID)) {
	if dualSampaID < 1024 {
		seg.bending.ForEachPad(f2cuid(padHandler))
	} else {
		seg.bending.ForEachPad(f2cuid(padHandler))
	}
}

func (seg *segmentation) GetNeighbours(paduid PadUID) []PadUID {
	cseg, p, err := seg.getCathSeg(paduid)
	if err != nil {
		return nil
	}
	padcid := cseg.GetNeighbours(p)
	var paduids []PadUID = make([]PadUID, len(padcid))
	offset := 0
	if !cseg.IsBending() {
		offset = seg.padUIDOffset
	}
	for i, _ := range padcid {
		paduids[i] = PadUID(int(padcid[i]) + offset)
	}
	return paduids
}

func (seg *segmentation) PadDualSampaChannel(paduid PadUID) int {
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

// ComputeSegmentationBBox returns the bounding box of
// the detection element represented by it segmentation.
func ComputeSegmentationBBox(seg Segmentation) geo.BBox {
	xmin := math.MaxFloat64
	ymin := xmin
	xmax := -xmin
	ymax := -ymin
	seg.ForEachPad(func(paduid PadUID) {
		bbox := ComputePadBBox(seg, paduid)
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

// ComputePadBBox returns the bounding box of one pad of the
// given segmentation.
func ComputePadBBox(seg Segmentation, paduid PadUID) geo.BBox {
	x := seg.PadPositionX(paduid)
	y := seg.PadPositionY(paduid)
	dx := seg.PadSizeX(paduid) / 2
	dy := seg.PadSizeY(paduid) / 2
	bbox, err := geo.NewBBox(x-dx, y-dy, x+dx, y+dy)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}
	return bbox
}
