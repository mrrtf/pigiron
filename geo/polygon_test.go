package geo

import (
	"testing"
)

var (
	testPolygon = Polygon{
		{0.1, 0.1},
		{1.1, 0.1},
		{1.1, 1.1},
		{2.1, 1.1},
		{2.1, 3.1},
		{1.1, 3.1},
		{1.1, 2.1},
		{0.1, 2.1},
		{0.1, 0.1}}

	testPolygon2 = Polygon{
		{-5.0, 10.0},
		{-5.0, -2.0},
		{0.0, -2.0},
		{0.0, -10.0},
		{5.0, -10.0},
		{5.0, 10.0},
		{-5.0, 10.0}}
)

func TestPolygonString(t *testing.T) {
	expected := "POLYGON (0.100000 0.100000,1.100000 0.100000,1.100000 1.100000,2.100000 1.100000,2.100000 3.100000,1.100000 3.100000,1.100000 2.100000,0.100000 2.100000,0.100000 0.100000)"
	if testPolygon.String() != expected {
		t.Errorf("expected string:%s and got:%s", expected, testPolygon.String())
	}
}
func TestCreateCounterClockwiseOrientedPolygon(t *testing.T) {
	p := Polygon{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
		{0, 0}}
	if !p.isCounterClockwiseOriented() {
		t.Error("polygon has wrong orientation")
	}
}

func TestCreateClockwiseOrientedPolygon(t *testing.T) {
	p := Polygon{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
		{0, 0}}
	if p.isCounterClockwiseOriented() {
		t.Error("polygon has wrong orientation")
	}
}

func TestSignedArea(t *testing.T) {
	expected := 4.0
	sa := testPolygon.signedArea()
	if !EqualFloat(sa, expected) {
		t.Errorf("expected signedArea to be %f but got %f", expected, sa)
	}
}

func TestAClosePolygonIsAPolygonWhereLastVertexIsTheSameAsTheFirstOne(t *testing.T) {
	if !testPolygon.isClosed() {
		t.Error("polygon should be closed")
	}
}

func TestClosingAClosedPolygonIsANop(t *testing.T) {
	closed, err := closePolygon(testPolygon)
	if err != nil {
		t.Error("closing failed")
	}
	if !EqualPolygon(testPolygon, closed) {
		t.Error("expected that closing an already closed polygon be a nop")
	}
}

func TestClosePolygon(t *testing.T) {
	opened := Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	expected := Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}
	closed, err := closePolygon(opened)
	if err != nil {
		t.Error("closing failed")
	}
	if !EqualPolygon(expected, closed) {
		t.Error("closing polygon yield unexpected result")
	}
}

func TestMustErrorIfClosingAPolygonResultInNonManhanttanPolygon(t *testing.T) {
	triangle := Polygon{{0, 0}, {1, 0}, {1, 1}}
	_, err := closePolygon(triangle)
	if err == nil {
		t.Error("closing should have yielded an error here")
	}
}

func TestAnOpenedPolygonCannotBeEqualToAClosedOneEventWithSameVertices(t *testing.T) {

	opened := Polygon{
		{0, 2}, {0, 0}, {2, 0}, {2, 4}, {1, 4}, {1, 2},
	}

	closed, _ := closePolygon(opened)

	if EqualPolygon(opened, closed) {
		t.Error("closed and opened polygon should not be equal")
	}
}

func TestPolygonsAreEqualAsLongAsTheyContainTheSameVertices(t *testing.T) {
	a := Polygon{{0, 2}, {0, 0}, {2, 0}, {2, 4}, {1, 4}, {1, 2}, {0, 2}}

	b := Polygon{{2, 4}, {2, 0}, {1, 4}, {1, 2}, {0, 2}, {0, 0}, {2, 4}}

	c := Polygon{{2, 4}, {2, 0}, {1, 4}, {1, 2}, {0, 2}, {0, 0}, {1, 3}}

	if !EqualPolygon(a, b) {
		t.Error("a and b polygons should be equal")
	}

	if EqualPolygon(b, c) {
		t.Error("b and c polygons should not be equal")
	}

	if EqualPolygon(Polygon{}, nil) {
		t.Error("empty polygon should be different from nil")
	}
}

func TestContainsMustErrorIfCalledOnNonClosedPolygon(t *testing.T) {
	opened := Polygon{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	_, err := opened.Contains(0, 0)
	if err == nil {
		t.Error("calling Contains on opened polygon should yield an error")
	}
}

func TestContainsReturnsTrueIfPointIsInsidePolygon(t *testing.T) {

	var testPoints = []struct {
		x, y float64
	}{
		{0, 0},
		{-4.999, -1.999},
	}
	for _, tp := range testPoints {
		ok, err := testPolygon2.Contains(tp.x, tp.y)
		if err != nil {
			t.Error("Contains should not trigger an error here")
		}
		if !ok {
			t.Errorf("testPolygon2 should contain (%f,%f)", tp.x, tp.y)
		}
	}
}

func TestContainsReturnsFalseIfPointIsExactlyOnAPolygonEdge(t *testing.T) {
	ok, err := testPolygon2.Contains(-2.5, -2)
	if err != nil {
		t.Error("Contains should not trigger an error here")
	}
	if ok {
		t.Error("testPolygon2 should not contain (-2.5,2)")
	}
}

func TestBBoxCreation(t *testing.T) {
	expected, _ := NewBBox(-5, -10, 5, 10)
	bbox := testPolygon2.BBox()
	if !EqualBBox(expected, bbox) {
		t.Errorf("expected bbox:%s and got:%s", expected.String(), bbox.String())
	}
}

func TestPolygonCenter(t *testing.T) {
	p := Polygon{{-80, -20}, {-70, -20}, {-70, -19.5}, {-80, -19.5}, {-80, -20}}
	b := p.BBox()
	if !EqualFloat(b.Xcenter(), -75.0) {
		t.Error("polygon xcenter is wrong")
	}
	if !EqualFloat(b.Ycenter(), -19.75) {
		t.Error("polygon ycenter is wrong")
	}
}

func squareDistanceTest(t *testing.T, p Polygon, point Vertex, expected float64) {
	d := SquaredDistancePointToPolygon(point, p)
	if !EqualFloat(d, expected) {
		t.Errorf("Expected distance from point (%f,%f) to polygon %s to be %f and got %f",
			point.X, point.Y, p.String(), expected, d)
	}
}
func TestPointOutsidePolygonDistanceToPolygonClosestToOneSegment(t *testing.T) {
	squareDistanceTest(t, testPolygon2, Vertex{-1, -6}, 1)
	squareDistanceTest(t, testPolygon2, Vertex{3, -14}, 16)
}
func TestPointOutsidePolygonDistanceToPolygonClosestToOneSegmentEndPoint(t *testing.T) {
	squareDistanceTest(t, testPolygon2, Vertex{-1, -14}, 17)
	squareDistanceTest(t, testPolygon2, Vertex{7, -14}, 20)
}

func TestPolygonTranslate(t *testing.T) {
	expected := Polygon{
		{-0.0, 0.0},
		{-0.0, -12.0},
		{5.0, -12.0},
		{5.0, -20.0},
		{10.0, -20.0},
		{10.0, 00.0},
		{0.0, 00.0}}

	tr := testPolygon2.Translate(5, -10)

	if !EqualPolygon(expected, tr) {
		t.Errorf("Translated polygon not as expected")
		t.Errorf("Want %s - Got %s", expected.String(), tr.String())
	}
}
