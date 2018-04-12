package mapping

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/aphecetche/pigiron/geo"
)

var ErrInvalidPadUID = errors.New("invalid pad uid")

const InvalidPadUID int = -1

type segmentation3 struct {
	segType                      int
	isBendingPlane               bool
	padGroups                    []padGroup
	padGroupTypes                []padGroupType
	padSizes                     []padSize
	dualSampaIDs                 []int
	padGroupIndex2PadUIDIndex    []int
	padUID2PadGroupTypeFastIndex []int
	padUID2PadGroupIndex         []int
}

func (seg *segmentation3) NofDualSampas() int {
	return len(seg.dualSampaIDs)
}

// NewSegmentation creates a segmentation object for the given detection element plane
func NewSegmentation(detElemID int, isBendingPlane bool) Segmentation {

	segType, err := detElemID2SegType(detElemID)
	if err != nil {
		return nil
	}
	builder := getSegmentationBuilder(segType)
	if builder == nil {
		return nil
	}
	return builder.Build(isBendingPlane)
}

func print(seg *segmentation3) {
	fmt.Println("segmentation3 has ", len(seg.dualSampaIDs), " dual sampa ids")
}

func newSegmentation(segType int, isBendingPlane bool, padGroups []padGroup,
	padGroupTypes []padGroupType, padSizes []padSize) *segmentation3 {

	seg := &segmentation3{segType, isBendingPlane, padGroups, padGroupTypes, padSizes, []int{}, []int{}, []int{}, []int{}}
	uniq := make(map[int]struct{})
	var empty struct{}
	for i := range padGroups {
		uniq[padGroups[i].fecID] = empty
	}
	for key := range uniq {
		seg.dualSampaIDs = append(seg.dualSampaIDs, key)
	}

	seg.init()
	return seg

}

func (seg *segmentation3) init() {
	paduid := 0

	for padGroupIndex := range seg.padGroups {
		seg.padGroupIndex2PadUIDIndex = append(seg.padGroupIndex2PadUIDIndex, paduid)
		pg := seg.padGroups[padGroupIndex]
		pgt := seg.padGroupTypes[pg.padGroupTypeID]
		for ix := 0; ix < pgt.nofPads; ix++ {
			for iy := 0; iy < pgt.nofPadsY; iy++ {
				if pgt.idByIndices(ix, iy) >= 0 {
					seg.padUID2PadGroupIndex = append(seg.padUID2PadGroupIndex, padGroupIndex)
					seg.padUID2PadGroupTypeFastIndex = append(seg.padUID2PadGroupTypeFastIndex, pgt.fastIndex(ix, iy))
					paduid++
				}
			}
		}
	}
}

func (seg *segmentation3) getPadUIDs(dualSampaID int) []int {
	pi := []int{}
	for pgi, pg := range seg.padGroups {
		if pg.fecID == dualSampaID {
			pgt := seg.padGroupTypes[pg.padGroupTypeID]
			i1 := seg.padGroupIndex2PadUIDIndex[pgi]
			for i := i1; i < i1+pgt.nofPads; i++ {
				pi = append(pi, i)
			}
		}
	}
	return pi
}

func (seg *segmentation3) DualSampaID(dualSampaIndex int) (int, error) {
	if dualSampaIndex >= len(seg.dualSampaIDs) {
		return -1, fmt.Errorf("Incorrect dualSampaIndex %d (should be within 0-%d range", dualSampaIndex,
			len(seg.dualSampaIDs))
	}
	return seg.dualSampaIDs[dualSampaIndex], nil
}

func (seg *segmentation3) NofPads() int {
	n := 0
	for i := 0; i < seg.NofDualSampas(); i++ {
		dsid, err := seg.DualSampaID(i)
		if err != nil {
			log.Fatalf("Could not get DualSampaID for i=%d", i)
		}
		n += len(seg.getPadUIDs(dsid))
	}
	return n
}

func (seg *segmentation3) ForEachPadInDualSampa(dualSampaID int, padHandler func(paduid int)) {
	for _, paduid := range seg.getPadUIDs(dualSampaID) {
		padHandler(paduid)
	}
}

func (seg *segmentation3) PadDualSampaChannel(paduid int) int {
	return seg.padGroupType(paduid).idByFastIndex(seg.padUID2PadGroupTypeFastIndex[paduid])
}

func (seg *segmentation3) PadDualSampaID(paduid int) int {
	return seg.padGroup(paduid).fecID
}

func (seg *segmentation3) padGroup(paduid int) padGroup {
	return seg.padGroups[seg.padUID2PadGroupIndex[paduid]]
}

func (seg *segmentation3) padGroupType(paduid int) padGroupType {
	return seg.padGroupTypes[seg.padGroup(paduid).padGroupTypeID]
}

func (seg *segmentation3) IsValid(paduid int) bool {
	return paduid != InvalidPadUID
}
func (seg *segmentation3) FindPadByFEE(dualSampaID, dualSampaChannel int) (int, error) {
	for _, paduid := range seg.getPadUIDs(dualSampaID) {
		if seg.padGroupType(paduid).idByFastIndex(seg.padUID2PadGroupTypeFastIndex[paduid]) == dualSampaChannel {
			return paduid, nil
		}
	}
	return InvalidPadUID, ErrInvalidPadUID
}

func (seg *segmentation3) FindPadByPosition(x, y float64) (int, error) {
	const epsilon = 1e-4
	xmin := x - epsilon
	xmax := x + epsilon
	ymin := y - epsilon
	ymax := y + epsilon
	box, err := geo.NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		return -1, fmt.Errorf("could not create bounding box: %v", err)
	}
	var pads []int
	for paduid := range seg.padGroups {
		pgrp := seg.padGroup(paduid)
		psiz := seg.padSizes[pgrp.padSizeID]
		x, y := seg.padPositionXY(paduid)
		xsiz := psiz.x * 0.5
		ysiz := psiz.y * 0.5
		b, err := geo.NewBBox(x-xsiz, y-ysiz, x+xsiz, y+ysiz)
		if err != nil {
			return -1, fmt.Errorf("could not create bounding box: %v", err)
		}
		_, err = geo.Intersect(b, box)
		if err == nil {
			pads = append(pads, paduid)
		}
	}

	dmin := +math.MaxFloat64
	id := InvalidPadUID
	for _, paduid := range pads {
		d := seg.squaredDistance(paduid, x, y)
		if d < dmin {
			dmin = d
			id = paduid
		}
	}

	if id == InvalidPadUID {
		return InvalidPadUID, fmt.Errorf("could not find pad")
	}

	return id, nil
}

func (seg *segmentation3) PadPositionX(paduid int) float64 {
	x, _ := seg.padPositionXY(paduid)
	return x
}

func (seg *segmentation3) PadPositionY(paduid int) float64 {
	_, y := seg.padPositionXY(paduid)
	return y
}

func (seg *segmentation3) padPositionXY(paduid int) (float64, float64) {
	pgrp := seg.padGroup(paduid)
	ptyp := seg.padGroupType(paduid)
	idx := seg.padUID2PadGroupTypeFastIndex[paduid]
	psiz := seg.padSizes[pgrp.padSizeID]
	x := pgrp.x + (float64(ptyp.ix(idx))+0.5)*psiz.x
	y := pgrp.y + (float64(ptyp.iy(idx))+0.5)*psiz.y
	return x, y
}

func (seg *segmentation3) PadSizeX(paduid int) float64 {
	pgrp := seg.padGroup(paduid)
	return seg.padSizes[pgrp.padSizeID].x
}

func (seg *segmentation3) PadSizeY(paduid int) float64 {
	pgrp := seg.padGroup(paduid)
	return seg.padSizes[pgrp.padSizeID].y
}

func (seg *segmentation3) squaredDistance(paduid int, x, y float64) float64 {
	px, py := seg.padPositionXY(paduid)
	px -= x
	py -= y
	return px*px + py*py
}
