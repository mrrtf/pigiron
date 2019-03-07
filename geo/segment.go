package geo

import (
	"fmt"
	"math"
)

type Segment struct {
	P0, P1 Vertex
}

// NewSegment creates a section going from
// (x1,y1) to (x2,y2)
func NewSegment(x1, y1, x2, y2 float64) *Segment {
	return &Segment{
		P0: Vertex{x1, y1},
		P1: Vertex{x2, y2},
	}
}

// NewVerticalSegment creates a vertical segment
// going from (x,y1) to (x,y2)
func NewVerticalSegment(x, y1, y2 float64) *Segment {
	return NewSegment(x, y1, x, y2)
}

// NewHorizontalSegment creates a horizontal segment
// going from (x1,y) to (x2,y)
func NewHorizontalSegment(y, x1, x2 float64) *Segment {
	return NewSegment(x1, y, x2, y)
}

// Vector returns a direction vector of this segment
func (s *Segment) Vector() Vertex {
	return VertexSub(s.P1, s.P0)
}

func (s Segment) String() string {
	msg := fmt.Sprintf("[ %v -> %v ]", s.P0, s.P1)
	if s.IsVertical() {
		msg += " V "
	}
	if s.IsHorizontal() {
		msg += " H "
	}
	return msg
}
func (s *Segment) Contains(v Vertex) bool {
	if !s.IsVertical() { // S is not  vertical
		if s.P0.X <= v.X && v.X <= s.P1.X {
			return true
		}
		if s.P0.X >= v.X && v.X >= s.P1.X {
			return true
		}
	} else { // S is vertical, so test Y  coordinate
		if s.P0.Y <= v.Y && v.Y <= s.P1.Y {
			return true
		}
		if s.P0.Y >= v.Y && v.Y >= s.P1.Y {
			return true
		}
	}
	return false
}

func (s *Segment) IsVertical() bool {
	// our float comparison is quite loose (1E-4)
	// but that's ok as we are supposed to deal
	// only with verticals and horizontals anyway
	return EqualFloat(s.P0.X, s.P1.X)
}

func (s *Segment) IsHorizontal() bool {
	// our float comparison is quite loose (1E-4)
	// but that's ok as we are supposed to deal
	// only with verticals and horizontals anyway
	return EqualFloat(s.P0.Y, s.P1.Y)
}

func AreSegmentParallel(s1 Segment, s2 Segment) bool {
	p := perp(s1.Vector(), s2.Vector())
	return math.Abs(p) < 1E-9
}

// IntersectSegmentLine returns the intersection point, if any,
// between the infinite line (defined here by a Segment) and
// the finite segment s.
func IntersectSegmentLine(s Segment, line Segment) (Vertex, bool) {
	if AreSegmentParallel(s, line) {
		return Vertex{}, false
	}

	v := line.Vector()
	u := s.Vector()
	w := VertexSub(s.P0, line.P0)

	si := -perp(v, w) / perp(v, u)

	x := s.P0.X + si*u.X
	y := s.P0.Y + si*u.Y

	return Vertex{x, y}, (si >= 0.0 && si <= 1.0)
}

// insersectTwoSegments finds the (2D) intersection of 2 segments.
// return (I0,I1,status)
// *I0 = intersect point (when it exists)
// *I1 =  endpoint of intersect segment [I0,I1] (when it exists)
// status = 0 disjoint (no intersect)
//            1=intersect  in unique point I0
//            2=overlap  in segment from I0 to I1
//
func IntersectTwoSegments(S1, S2 Segment) (Vertex, Vertex, int) {
	u := VertexSub(S1.P1, S1.P0) // Vector    u = S1.P1 - S1.P0;
	v := VertexSub(S2.P1, S2.P0) // Vector    v = S2.P1 - S2.P0;
	w := VertexSub(S1.P0, S2.P0) // Vector    w = S1.P0 - S2.P0;

	D := perp(u, v) //float     D = perp(u,v);

	const SMALL_NUM float64 = 1E-9

	// test if  they are parallel (includes either being a point)
	if math.Abs(D) < SMALL_NUM { // S1 and S2 are parallel
		if perp(u, w) != 0 || perp(v, w) != 0 {
			return Vertex{}, Vertex{}, 0 // they are NOT collinear
		}
		// they are collinear or degenerate
		// check if they are degenerate  points
		du := dot(u, u)
		dv := dot(v, v)
		if du == 0 && dv == 0 { // both segments are points
			if !EqualVertex(S1.P0, S2.P0) {
				// they are distinct  points
				return Vertex{}, Vertex{}, 0
			}
			return S1.P0, Vertex{}, 1 // they are the same point
		}
		if du == 0 { // S1 is a single point
			if !S2.Contains(S1.P0) {
				// but is not in S2
				return Vertex{}, Vertex{}, 0
			}
			return S1.P0, Vertex{}, 1
		}
		if dv == 0 { // S2 a single point
			if !S1.Contains(S2.P0) {
				// but is not in S1
				return Vertex{}, Vertex{}, 0
			}
			return S2.P0, Vertex{}, 1
		}
		// they are collinear segments - get  overlap (or not)
		var t0, t1 float64            // endpoints of S1 in eqn for S2
		w2 := VertexSub(S1.P1, S2.P0) //Vector w2 = S1.P1 - S2.P0;
		if v.X != 0 {
			t0 = w.X / v.X
			t1 = w2.X / v.X
		} else {
			t0 = w.Y / v.Y
			t1 = w2.Y / v.Y
		}
		if t0 > t1 { // must have t0 smaller than t1
			// swap if not
			t0, t1 = t1, t0
		}
		if t0 > 1 || t1 < 0 {
			return Vertex{}, Vertex{}, 0 // NO overlap
		}
		if t0 < 0 {
			t0 = 0 // clip to min 0
		}
		if t1 > 1 {
			t1 = 1 // clip to max 1
		}
		if EqualFloat(t0, t1) {
			// intersect is a point
			return VertexAdd(S2.P0, v.Scale(t0)), Vertex{}, 1
		}

		// they overlap in a valid subsegment
		return VertexAdd(S2.P0, v.Scale(t0)), VertexAdd(S2.P0, v.Scale(t1)), 2
	}

	// the segments are skew and may intersect in a point
	// get the intersect parameter for S1
	sI := perp(v, w) / D
	if sI < 0 || sI > 1 {
		// no intersect with S1
		return Vertex{}, Vertex{}, 0
	}

	// get the intersect parameter for S2
	tI := perp(u, w) / D
	if tI < 0 || tI > 1 {
		// no intersect with S2
		return Vertex{}, Vertex{}, 0
	}

	// compute S1 intersect point
	return VertexAdd(S1.P0, u.Scale(sI)), Vertex{}, 1
}
