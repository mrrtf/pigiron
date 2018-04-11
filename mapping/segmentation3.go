package mapping

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/aphecetche/pigiron/geo"
)

// ErrInvalidPadUID signals that an invalid pad uid was used / returned
var ErrInvalidPadUID = errors.New("invalid pad uid")

// InvalidPadUID is a special integer that denotes a non existing pad
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
	grid                         padGroupGrid
	detElemID                    int
}

func (seg *segmentation3) DetElemID() int {
	return seg.detElemID
}

func (seg *segmentation3) setDetElemID(de int) {
	seg.detElemID = de
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
	seg := builder.Build(isBendingPlane)
	seg.setDetElemID(detElemID)
	return seg
}

func (seg *segmentation3) Print(out io.Writer) {
	fmt.Fprintf(out, "segmentation3 has %v dual sampa ids = %v\n", len(seg.dualSampaIDs), seg.dualSampaIDs)
	d := make(map[int]int)
	for i := 0; i < seg.NofDualSampas(); i++ {
		dsid, err := seg.DualSampaID(i)
		if err != nil {
			panic(err)
		}
		d[dsid] = 1
	}
	fmt.Fprintf(out, "cells=%v\n", seg.grid.cells)
	for c := range seg.grid.cells {
		fmt.Fprintf(out, "%v ", seg.padGroups[c].fecID)
	}
	fmt.Fprintf(out, "\n")
	seg.grid.Print(out)
}

func newSegmentation(segType int, isBendingPlane bool, padGroups []padGroup,
	padGroupTypes []padGroupType, padSizes []padSize) *segmentation3 {
	seg := &segmentation3{
		segType:        segType,
		isBendingPlane: isBendingPlane,
		padGroups:      padGroups,
		padGroupTypes:  padGroupTypes,
		padSizes:       padSizes}
	uniq := make(map[int]struct{})
	var empty struct{}
	for i := range padGroups {
		uniq[padGroups[i].fecID] = empty
	}
	for key := range uniq {
		seg.dualSampaIDs = append(seg.dualSampaIDs, key)
	}
	seg.init()
	//seg.Print(os.Stdout)
	return seg
}

func (seg *segmentation3) init() {
	// here must make two loops
	// first one to fill in the 3 index slices
	// - padGroupIndex2PadUIDIndex
	// - padUID2PadGroupIndex
	// - padUID2PadGroupTypeFastIndex
	// then compute the global x,y ranges to be able to compute a grid
	// then loop over to put each cell (a cell is a pair (box,paduid))
	// within the correct grid cellSlice

	seg.fillIndexSlices()
	bbox := seg.computeBbox()
	seg.fillGrid(bbox)
}

func (seg *segmentation3) computeBbox() geo.BBox {
	xmin := math.MaxFloat64
	ymin := xmin
	xmax := -xmin
	ymax := -ymin
	for paduid := range seg.padUID2PadGroupTypeFastIndex {
		x := seg.PadPositionX(paduid)
		y := seg.PadPositionY(paduid)
		dx := seg.PadSizeX(paduid) / 2
		dy := seg.PadSizeY(paduid) / 2
		xmin = math.Min(xmin, x-dx)
		xmax = math.Max(xmax, x+dx)
		ymin = math.Min(ymin, y-dy)
		ymax = math.Max(ymax, y+dy)
	}

	bbox, err := geo.NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		panic(err)
	}
	return bbox
}

func (seg *segmentation3) fillIndexSlices() {
	paduid := 0
	for padGroupIndex := range seg.padGroups {
		seg.padGroupIndex2PadUIDIndex = append(seg.padGroupIndex2PadUIDIndex, paduid)
		pg := seg.padGroups[padGroupIndex]
		pgt := seg.padGroupTypes[pg.padGroupTypeID]
		for ix := 0; ix < pgt.NofPadsX; ix++ {
			for iy := 0; iy < pgt.NofPadsY; iy++ {
				if pgt.idByIndices(ix, iy) >= 0 {
					seg.padUID2PadGroupIndex = append(seg.padUID2PadGroupIndex, padGroupIndex)
					seg.padUID2PadGroupTypeFastIndex = append(seg.padUID2PadGroupTypeFastIndex, pgt.fastIndex(ix, iy))
					paduid++
				}
			}
		}
	}
}

func (seg *segmentation3) padGroupBox(i int) geo.BBox {
	pg := seg.padGroups[i]
	pgt := seg.padGroupTypes[pg.padGroupTypeID]
	dx := seg.padSizes[pg.padSizeID].x * float64(pgt.NofPadsX)
	dy := seg.padSizes[pg.padSizeID].y * float64(pgt.NofPadsY)
	x := pg.x
	y := pg.y
	box, err := geo.NewBBox(x, y, x+dx, y+dy)
	if err != nil {
		panic(err)
	}
	return box
}

func (seg *segmentation3) fillGrid(bbox geo.BBox) {
	// for each cell in the grid we find out which
	// padgroups have their bounding box intersecting with
	// the cell bounding box and insert them in that cell
	// if the intersect is not nil

	seg.grid = *(newPadGroupGrid(bbox))
	for index := range seg.grid.cells {
		cbox := seg.grid.cellBox(index)
		for i := range seg.padGroups {
			pbox := seg.padGroupBox(i)
			_, err := geo.Intersect(cbox, pbox)
			if err == nil {
				seg.grid.cells[index] = append(seg.grid.cells[index], i)
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
			for i := i1; i < i1+pgt.NofPads; i++ {
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

func (seg *segmentation3) PadSizeX(paduid int) float64 {
	return seg.padSizes[seg.padGroup(paduid).padSizeID].x
}
func (seg *segmentation3) PadSizeY(paduid int) float64 {
	return seg.padSizes[seg.padGroup(paduid).padSizeID].y
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

func (seg *segmentation3) padGroup(paduid int) *padGroup {
	return &seg.padGroups[seg.padUID2PadGroupIndex[paduid]]
}

func (seg *segmentation3) padGroupType(paduid int) *padGroupType {
	return &seg.padGroupTypes[seg.padGroup(paduid).padGroupTypeID]
}

func (seg *segmentation3) FindPadByPosition(x, y float64) (int, error) {
	pgis := seg.grid.padGroupIndex(x, y)
	var pgti []int
	for pgi := range pgis {
		pgIndex := pgis[pgi]
		pg := seg.padGroups[pgIndex]
		pgt := seg.padGroupTypes[pg.padGroupTypeID]
		lx := x - pg.x
		ly := y - pg.y
		ix := int(math.Trunc(lx / seg.padSizes[pg.padSizeID].x))
		iy := int(math.Trunc(ly / seg.padSizes[pg.padSizeID].y))
		valid := pgt.areIndicesPossible(ix, iy) && lx >= 0.0 && ly >= 0.0
		if valid {
			// find in padUID2PadGroupTypeFastIndex array, starting at seg.padGroupIndex2PadUIDIndex[pgis[pgi]]
			// the paduid corresponding to pgt.fastIndex(ix,iy)
			// FIXME : that is wrong.
			a := seg.padGroupIndex2PadUIDIndex[pgIndex]
			asize := len(seg.padGroupIndex2PadUIDIndex) - 1
			var b int
			if pgIndex >= asize-1 {
				b = seg.NofPads() - 1
			} else {
				b = seg.padGroupIndex2PadUIDIndex[pgIndex+1]
			}
			for j := a; j <= b; j++ {
				if pgt.fastIndex(ix, iy) == seg.padUID2PadGroupTypeFastIndex[j] {
					pgti = append(pgti, j)
					break
				}
			}
		}
	}
	if len(pgti) > 1 {
		log.Fatalf("more than one match ! %v", pgti)
	}
	if len(pgti) > 0 {
		return pgti[0], nil
	}
	return InvalidPadUID, ErrInvalidPadUID
}

func (seg *segmentation3) PadPositionX(paduid int) float64 {
	pg := seg.padGroup(paduid)
	pgt := seg.padGroupType(paduid)
	return pg.x + (float64(pgt.ix(seg.padUID2PadGroupTypeFastIndex[paduid]))+0.5)*
		seg.padSizes[pg.padSizeID].x
}

func (seg *segmentation3) PadPositionY(paduid int) float64 {
	pg := seg.padGroup(paduid)
	pgt := seg.padGroupType(paduid)
	return pg.y + (float64(pgt.iy(seg.padUID2PadGroupTypeFastIndex[paduid]))+0.5)*
		seg.padSizes[pg.padSizeID].y
}
