package geo

import (
	"testing"
)

func TestClipperVertexInside(t *testing.T) {

	// From Fig. 6 example from "Reetrant Polygon Clipping",
	// Ivan E. Sutherland and Gary W. Hodgman
	// Communications of the ACM, January 1974, Volume 17, Number 1, p32-42

	p2 := Vertex{20.0, 20.0}

	right := Segment{Vertex{10, -10}, Vertex{10, 10}}
	left := Segment{Vertex{-10, 10}, Vertex{-10, -10}}
	top := Segment{Vertex{10, 10}, Vertex{-10, 10}}
	bottom := Segment{Vertex{-10, -10}, Vertex{10, -10}}

	if isInside(p2, right) != false {
		t.Errorf("Expected p2 to be hidden from right edge")
	}
	if isInside(p2, left) != true {
		t.Errorf("Expected p2 to be visible from left edge")
	}
	if isInside(p2, top) != false {
		t.Errorf("Expected p2 to be hidden from top edge")
	}
	if isInside(p2, bottom) != true {
		t.Errorf("Expected p2 to be visible from top edge")
	}
}

func TestClipPolygon(t *testing.T) {

	// From Fig. 6 example from "Reetrant Polygon Clipping",
	// Ivan E. Sutherland and Gary W. Hodgman
	// Communications of the ACM, January 1974, Volume 17, Number 1, p32-42

	input := Polygon{
		{-20.0, 20.0},
		{0.0, -20.0},
		{20.0, -20.0},
		{20.0, 20.0},
		{-20.0, 20.0}}

	window, _ := NewBBox(-10.0, -10.0, 10.0, 10.0)

	c, err := ClipPolygon(input, window)

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := Polygon{
		{-10.0, 10.0},
		{-10.0, 0.0},
		{-5, -10},
		{10, -10},
		{10, 10},
		{-10, 10},
	}

	if !EqualPolygon(expected, c) {
		t.Errorf("Want polygon=%v\n.Got %v\n", expected, c)

	}
}
