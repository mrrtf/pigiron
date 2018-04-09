package geo

import "testing"

var dummy float64

func TestVerticalLeftEdgeIsTopToBottom(t *testing.T) {
	edge := verticalEdge{dummy, 12, 1}
	if !isLeftEdge(edge) || !isTopToBottom(edge) {
		t.Error("edge should be a left edge, aka top to bottom")
	}
}

func TestVerticalRightEdgeIsBottomToTop(t *testing.T) {
	edge := verticalEdge{dummy, 1, 12}
	if !isRightEdge(edge) || !isBottomToTop(edge) {
		t.Error("edge should be a right edge, aka bottom to top")
	}
}

func TestLeftToRightHorizontalEdgeHasEndPointGreaterThanStartPoint(t *testing.T) {
	edge := horizontalEdge{dummy, 1, 12}
	if !isLeftToRight(edge) {
		t.Error("edge should be left to right")
	}
}

func TestRightToLeftHorizontalEdgeHasEndPointSmallerThanStartPoint(t *testing.T) {
	edge := horizontalEdge{dummy, 12, 1}
	if !isRightToLeft(edge) {
		t.Error("edge should be right to left")
	}
}

func TestVerticalEdgeWithBeginAboveEndIsALefty(t *testing.T) {
	v := verticalEdge{0, 12, 10}
	if !isLeftEdge(v) || isRightEdge(v) {
		t.Error("edge should be a left edge")
	}
}

func TestVerticalEdgeWithBeginAboveEndIsARighty(t *testing.T) {
	v := verticalEdge{0, 10, 12}
	if !isRightEdge(v) || isLeftEdge(v) {
		t.Error("edge should be a right edge")
	}
}

func TestVerticalEdgeHasTopAndBottom(t *testing.T) {
	v := verticalEdge{2, 10, 12}
	if !EqualFloat(bottom(v), 10) || !EqualFloat(top(v), 12) {
		t.Error("something is wrong with top and bottom")
	}
}

func TestBeginAndEndForLeftVerticalEdge(t *testing.T) {
	v := verticalEdge{0, 7, 1}
	if !EqualFloat(top(v), 7) ||
		!EqualFloat(bottom(v), 1) ||
		!EqualVertex(v.begin(), Vertex{0, 7}) ||
		!EqualVertex(v.end(), Vertex{0, 1}) {
		t.Error("something is wrong with begin and end for left vertical edge")
	}
}

func TestBeginAndEndForRightVerticalEdge(t *testing.T) {
	v := verticalEdge{0, 1, 7}
	if !EqualFloat(top(v), 7) ||
		!EqualFloat(bottom(v), 1) ||
		!EqualVertex(v.begin(), Vertex{0, 1}) ||
		!EqualVertex(v.end(), Vertex{0, 7}) {
		t.Error("something is wrong with begin and end for right vertical edge")
	}
}

func TestBeginAndEndForALeftToRightHorizontalEdge(t *testing.T) {
	h := horizontalEdge{0, 1, 7}
	if !EqualFloat(left(h), 1) ||
		!EqualFloat(right(h), 7) ||
		!EqualVertex(h.begin(), Vertex{1, 0}) ||
		!EqualVertex(h.end(), Vertex{7, 0}) {
		t.Error("something is wrong with begin and end for left to right horizontal edge")
	}
}

func TestBeginAndEndForARightToLeftHorizontalEdge(t *testing.T) {
	h := horizontalEdge{0, 7, 1}
	if !EqualFloat(left(h), 1) ||
		!EqualFloat(right(h), 7) ||
		!EqualVertex(h.begin(), Vertex{7, 0}) ||
		!EqualVertex(h.end(), Vertex{1, 0}) {
		t.Error("something is wrong with begin and end for right to left horizontal edge")
	}
}
