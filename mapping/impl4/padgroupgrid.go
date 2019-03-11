package impl4

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/mrrtf/pigiron/geo"
)

type padGroupGrid struct {
	bbox   geo.BBox
	nx, ny int
	cells  [][]int // grid[ix+iy*ngy] = []int{ list of padgroupindices in this grid cell}
	gx     float64
	gy     float64
}

func (g *padGroupGrid) cellBox(index int) geo.BBox {
	ix, iy := g.getIndices(index)
	x := g.bbox.Xmin() + float64(ix)*g.gx
	y := g.bbox.Ymin() + float64(iy)*g.gy
	box, err := geo.NewBBox(x, y, x+g.gx, y+g.gy)
	if err != nil {
		panic(err)
	}
	return box
}

func newPadGroupGrid(bbox geo.BBox) *padGroupGrid {
	gx, gy := 10.0, 10.0
	eps := 1E-4
	ix, iy := getGridIndices(bbox.Width()-eps, bbox.Height()-eps, gx, gy)
	pgg := &padGroupGrid{
		bbox:  bbox,
		nx:    ix + 1,
		ny:    iy + 1,
		cells: make([][]int, (ix+1)*(iy+1)),
		gx:    gx,
		gy:    gy,
	}
	return pgg
}

func getGridIndices(x, y, cellWidth, cellHeight float64) (int, int) {
	if x < 0.0 || y < 0.0 {
		panic(fmt.Sprintf("x,y are not >=0 : x=%v y=%v", x, y))
	}
	return int(math.Trunc(x / cellWidth)), int(math.Trunc(y / cellHeight))
}

var (
	errCoordOutOfRange = errors.New("(x,y) coordinates outside of segmentation bounding box")
	invalidIndex       = -1
)

func (g *padGroupGrid) getGridIndices(x, y float64) (int, int, error) {
	if !g.bbox.Contains(x, y) {
		return invalidIndex, invalidIndex, errCoordOutOfRange
	}
	ix, iy := getGridIndices(x-g.bbox.Xmin(), y-g.bbox.Ymin(), g.gx, g.gy)
	return ix, iy, nil
}

func (g *padGroupGrid) getIndices(index int) (int, int) {
	//index := ix + iy*g.nx
	iy := index / g.nx
	ix := index - iy*g.nx
	return ix, iy
}

func (g *padGroupGrid) Print(out io.Writer) {
	fmt.Fprintf(out, "padGroupGrid: nx %v ny %v %v bins bbox %v\n", g.nx, g.ny, len(g.cells), g.bbox)
	for i, c := range g.cells {
		ix, iy := g.getIndices(i)
		fmt.Fprintf(out, "padGroupGrid: cell %3d (%2d,%2d) has %5d elements\n", i, ix, iy, len(c))
	}
}

func (g *padGroupGrid) getIndex(x, y float64) int {
	ix, iy, err := g.getGridIndices(x, y)
	if err != nil {
		return invalidIndex
	}
	i := ix + iy*g.nx
	if i < 0 {
		panic(fmt.Sprintf("x %v y %v grid=%v", x, y, g))
	}
	if i >= len(g.cells) {
		return invalidIndex
	}

	return i
}

func (g *padGroupGrid) padGroupIndex(x, y float64) []int {
	i := g.getIndex(x, y)
	if i == invalidIndex {
		return nil
	}
	if i < 0 || i >= len(g.cells) {
		log.Fatalf("i out of bounds %d vs %d", i, len(g.cells))
	}
	return g.cells[i]
}
