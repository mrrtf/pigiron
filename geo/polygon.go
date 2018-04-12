package geo

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// Polygon describes a simple, rectilinear, closed set of vertices
// with a specific orientation
type Polygon []Vertex

func (p Polygon) isManhattan() bool {
	for i := 0; i < len(p)-1; i++ {
		if !isVerticalSegment(p[i], p[i+1]) &&
			!isHorizontalSegment(p[i], p[i+1]) {
			return false
		}
	}
	return true
}

func (p Polygon) isCounterClockwiseOriented() bool {
	return p.signedArea() > 0
}

func (p Polygon) signedArea() float64 {
	/// Compute the signed area of this polygon
	/// Algorithm from F. Feito, J.C. Torres and A. Urena,
	/// Comput. & Graphics, Vol. 19, pp. 595-600, 1995
	area := 0.0
	for i := 0; i < len(p)-1; i++ {
		current := p[i]
		next := p[i+1]
		area += current.X*next.Y - next.X*current.Y
	}
	return area * 0.5
}

func (p Polygon) isClosed() bool {
	return p[0] == p[len(p)-1]
}

// EqualPolygon checks if two polygon are the same
// (same vertices, whetever the order)
func EqualPolygon(a, b Polygon) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	sa := a.getSortedVertices()
	sb := b.getSortedVertices()

	for i, v := range sa {
		if !EqualVertex(v, sb[i]) {
			return false
		}
	}
	return true
}

func closePolygon(p Polygon) (Polygon, error) {
	if p.isClosed() {
		return p, nil
	}
	np := Polygon{}
	np = append(np, p...)
	np = append(np, p[0])
	if !np.isManhattan() {
		return nil, errors.New("closing resulted in non Manhattan polygon")
	}
	return np, nil
}

func (p Polygon) String() string {
	s := fmt.Sprintf("POLYGON (")
	for i, v := range p {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%f %f", v.X, v.Y)
	}
	s += ")"
	return s
}

func (p Polygon) getVertices() []Vertex {
	size := len(p)
	if p.isClosed() {
		size--
	}
	c := make([]Vertex, size)
	copy(c, p)
	return c
}

func sortVertices(vertices []Vertex) []Vertex {
	c := append([]Vertex{}, vertices...)
	sort.Slice(c, func(i, j int) bool {
		if EqualFloat(c[i].X, c[j].X) {
			return c[i].Y < c[j].Y
		}
		return c[i].X < c[j].X
	})
	return c
}

func (p Polygon) getSortedVertices() []Vertex {
	return sortVertices(p.getVertices())
}

// Contains returns true if the point (xp,yp) is inside the polygon
//
// Note that this algorithm yields unpredicatable result if the point xp,yp
// is on one edge of the polygon. Should not generally matters, except when comparing
// two different implementations maybe.
//
// TODO : look e.g. to http://alienryderflex.com/polygon/ for some possible optimizations
// (e.g. pre-computation)
func (p Polygon) Contains(xp, yp float64) (bool, error) {
	if !p.isClosed() {
		return false, errors.New("Contains can only work with closed polygons")
	}

	pj := p[len(p)-1]
	oddNodes := false
	for _, pi := range p {
		if (pi.Y < yp && pj.Y >= yp) || (pj.Y < yp && pi.Y >= yp) {
			if pi.X+(yp-pi.Y)/(pj.Y-pi.Y)*(pj.X-pi.X) < xp {
				oddNodes = !oddNodes
			}
		}
		pj = pi
	}
	return oddNodes, nil
}

// BBox returns the bounding box of the polygon.
func (p Polygon) BBox() BBox {
	return getVerticesBBox(p.getVertices())
}

// SquaredDistancePointToPolygon return the square of the distance
// between a point and a polygon
func SquaredDistancePointToPolygon(point Vertex, p Polygon) float64 {
	d := math.MaxFloat64
	for i := 0; i < len(p)-1; i++ {
		s0 := p[i]
		s1 := p[i+1]
		d2 := SquaredDistanceOfPointToSegment(point, s0, s1)
		d = math.Min(d, d2)
	}
	return d
}

func areCounterClockwisePolygons(polygons []Polygon) bool {
	for _, p := range polygons {
		if !p.isCounterClockwiseOriented() {
			return false
		}
	}
	return true
}
