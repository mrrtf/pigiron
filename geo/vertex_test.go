package geo

import "testing"

var (
	p0 = Vertex{0, 0}
	p1 = Vertex{6, 0}
)

func TestVertical(t *testing.T) {
	v1 := Vertex{12, 0}
	v2 := Vertex{12, 20}
	if !isVerticalSegment(v1, v2) {
		t.Error("segment v1-v2 should be vertical")
	}
	if isVerticalSegment(v1, p0) {
		t.Error("segment v1-p0 should not be vertical")
	}
}
func TestHorizontal(t *testing.T) {
	v1 := Vertex{0, 12}
	v2 := Vertex{20, 12}
	if !isHorizontalSegment(v1, v2) {
		t.Error("segment v1-v2 should be horizontal")
	}
	if isHorizontalSegment(v1, p0) {
		t.Error("segment v1-p0 should not be horizontal")
	}
}

func TestVertexEquality(t *testing.T) {
	a := Vertex{0, 1}
	b := Vertex{0, 1.0 + 1E-6}
	if !EqualVertex(a, b) {
		t.Errorf("Vertices %s and %s should be equal", a.String(), b.String())
	}
}

func checkD2P(t *testing.T, p Vertex, expected float64) {
	d := SquaredDistanceOfPointToSegment(p, p0, p1)
	if !EqualFloat(d, expected) {
		t.Errorf("expected distance to be %e and got %e", expected, d)
	}
}

func TestDistancePointToSegmentWhereBasePointIsWithinSegment(t *testing.T) {
	checkD2P(t, Vertex{1.5, 3.5}, 12.25)
}

func TestDistancePointToSegmentWhereBasePointIfLeftOfSegment(t *testing.T) {
	checkD2P(t, Vertex{-3, 3}, 18)
}

func TestDistancePointToSegmentWhereBasePointIsRightOfSegment(t *testing.T) {
	checkD2P(t, Vertex{8, 2}, 8)
}

func TestVertexString(t *testing.T) {
	a := Vertex{12.34, -56.789}
	s := a.String()
	expected := "( 12.34 -56.789 )"
	if s != expected {
		t.Errorf("expected string:%s and got:%s", expected, s)
	}
}
