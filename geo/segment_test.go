package geo_test

import (
	"testing"

	"github.com/aphecetche/pigiron/geo"
)

var (
	ymin, ymax, x0, x2 float64
	y3, x3min, x3max   float64
	s1, s2, s3, s4     *geo.Segment
	x4min, x4max       float64
)

func init() {
	ymin = -1.0
	ymax = 1.0

	x0 = 0.0
	x2 = 2.0
	s1 = geo.NewVerticalSegment(x0, ymin, ymax)
	s2 = geo.NewVerticalSegment(x2, ymin, ymax)

	y3 = 0.25
	x3min = -1.0
	x3max = 1.0
	s3 = geo.NewHorizontalSegment(y3, x3min, x3max)

	x4min = -0.12
	x4max = 1.24
	s4 = geo.NewHorizontalSegment(y3, x4min, x4max)
}

func TestSegmentIntersectTwoSegmentThatDoNotIntersect(t *testing.T) {
	i0, i1, status := geo.IntersectTwoSegments(*s1, *s2)
	if status != 0 {
		t.Errorf("Want status=0. Got %d\n", status)
	}
	if !geo.EqualVertex(geo.Vertex{}, i0) {
		t.Errorf("Want default i0. Got %v\n", i0)
	}
	if !geo.EqualVertex(geo.Vertex{}, i1) {
		t.Errorf("Want default i1. Got %v\n", i1)
	}
}

func TestSegmentIntersectTwoPerpSegmentThatDoIntersect(t *testing.T) {
	i0, i1, status := geo.IntersectTwoSegments(*s1, *s3)
	if status != 1 {
		t.Errorf("Want status=1. Got %d\n", status)
	}
	if !geo.EqualVertex(geo.Vertex{}, i1) {
		t.Errorf("Want default i1. Got %v\n", i1)
	}
	expected := geo.Vertex{X: x0, Y: y3}
	if !geo.EqualVertex(i0, expected) {
		t.Errorf("Want vertex=%v. Got %v\n", expected, i0)
	}
}

func TestSegmentIntersectTwoOverlappingSegments(t *testing.T) {
	i0, i1, status := geo.IntersectTwoSegments(*s3, *s4)
	if status != 2 {
		t.Errorf("Want status=2. Got %d\n", status)
	}
	e0 := geo.Vertex{X: x4min, Y: y3}
	if !geo.EqualVertex(e0, i0) {
		t.Errorf("Want vertex=%v. Got %v\n", e0, i0)
	}
	e1 := geo.Vertex{X: x3max, Y: y3}
	if !geo.EqualVertex(e1, i1) {
		t.Errorf("Want vertex=%v. Got %v\n", e1, e1)
	}
}

func TestSegmentParallel(t *testing.T) {

	s1 := geo.Segment{P0: geo.Vertex{X: 0.0, Y: 0.0}, P1: geo.Vertex{X: 1.0, Y: 1.0}}
	s2 := geo.Segment{P0: geo.Vertex{X: 1.0, Y: 0.0}, P1: geo.Vertex{X: 2.0, Y: 1.0}}
	s3 := geo.Segment{P0: geo.Vertex{X: 1.0, Y: 0.0}, P1: geo.Vertex{X: 1.0, Y: 1.0}}

	if !geo.AreSegmentParallel(s1, s2) {
		t.Errorf("Segments should be parallel")
	}

	if geo.AreSegmentParallel(s1, s3) {
		t.Errorf("Segments should NOT be parallel")
	}

	se := geo.Segment{P0: geo.Vertex{X: 20, Y: -20}, P1: geo.Vertex{X: 20, Y: 20}}
	clipEdge := geo.Segment{P0: geo.Vertex{X: 10, Y: 10}, P1: geo.Vertex{X: -10, Y: 10}}
	if geo.AreSegmentParallel(se, clipEdge) {
		t.Errorf("Segments se and clipEdge should NOT be parallel")
	}
}

func TestSegmentIntersectSegmentLine(t *testing.T) {

	s := geo.Segment{P0: geo.Vertex{X: -20, Y: 20}, P1: geo.Vertex{X: 0, Y: -20}}
	line := geo.Segment{P0: geo.Vertex{X: -10, Y: -10}, P1: geo.Vertex{X: 10, Y: -10}}

	i, ok := geo.IntersectSegmentLine(s, line)

	if !ok {
		t.Errorf("Should have an intersection here")
	}

	expected := geo.Vertex{X: -5, Y: -10}
	if !geo.EqualVertex(i, expected) {
		t.Errorf("Want %v. Got %v\n", expected, i)
	}

	s = geo.Segment{P0: geo.Vertex{X: -20, Y: 20}, P1: geo.Vertex{X: 10, Y: 1}}
	_, ok = geo.IntersectSegmentLine(s, line)
	if ok {
		t.Errorf("Should not get an intersect here")
	}

}
