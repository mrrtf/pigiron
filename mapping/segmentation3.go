package mapping

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/aphecetche/pigiron/geo"
)

// ErrInvalidPadUID signals that an invalid pad uid was used / returned
var ErrInvalidPadUID = errors.New("invalid pad uid")

// InvalidPadUID is a special integer that denotes a non existing pad
const InvalidPadUID int = -1

type padGroupGrid struct {
	bbox   geo.BBox
	nx, ny int
	grid   [][]int // grid[ix+iy*ngy] = []int{ list of padgroupindices in this grid cell}
}
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
}

func newPadGroupGrid(bbox geo.BBox) *padGroupGrid {
	nx, ny := getGridIndices(bbox.Width(), bbox.Height())
	pgg := &padGroupGrid{
		bbox: bbox,
		nx:   nx,
		ny:   ny,
		grid: make([][]int, nx*ny),
	}
	return pgg
}

func getGridIndices(x, y float64) (int, int) {
	gx := 20.0 // cm
	gy := 20.0 // cm
	return int(x / gx), int(y / gy)
}

func (g *padGroupGrid) String() string {
	return fmt.Sprintf("padGroupGrid: %v bins", len(g.grid))
}

func (g *padGroupGrid) getIndex(x, y float64) int {
	ix, iy := getGridIndices(x, y)
	return ix + iy*g.nx
}

func (g *padGroupGrid) padGroupIndex(x, y float64) []int {
	return g.grid[g.getIndex(x, y)]
}

func (g *padGroupGrid) insertPadGroup(x, y float64, padGroupIndex int) {
	g.grid[g.getIndex(x, y)] = append(g.grid[g.getIndex(x, y)], padGroupIndex)
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
	return seg
}

func (seg *segmentation3) report() {
	fmt.Println(&seg.grid)
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

	seg.report()
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
	fmt.Println(bbox.String(), " npads=", len(seg.padUID2PadGroupIndex))
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

func (seg *segmentation3) fillGrid(bbox geo.BBox) {
	/*
		box, err := getPadBox(ix, iy, pg.x, pg.y, dx, dy)
		if err != nil {
			panic(err)
		}
		fmt.Printf("box:%s\n", box.String())
			gx := int(box.Xcenter() / seg.gridSizeX)
			gy := int(box.Ycenter() / seg.gridSizeY)
			if seg.grid == nil {
				seg.grid = [][]cellSlice{}
			}
			if seg.grid[gx] == nil {
				seg.grid[gx] = []cellSlice{}
			}
			if seg.grid[gx][gy] == nil {
				seg.grid[gx][gy] = cellSlice{}
			}
			fmt.Printf("padGroupIndex %3d gx %3d gy %3d", padGroupIndex, gx, gy)
			seg.grid[gx][gy] = append(seg.grid[gx][gy], cell{box, paduid})
	*/
}

func (seg *segmentation3) getPadUIDs(dualSampaID int) []int {
	pi := []int{}
	for pgi := range seg.padGroups {
		pg := seg.padGroups[pgi]
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

/// FIXME : to be implemented...
func (seg *segmentation3) FindPadByPosition(x, y float64) (int, error) {
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
