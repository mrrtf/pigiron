package geo

import (
	"fmt"
	"testing"
)

func TestNewContourGeneratesEmptyContourForEmptyInput(t *testing.T) {
	polygons := []Polygon{}
	c, err := NewContour(polygons)
	if err != nil {
		t.Error("should not trigger an error here")
	}
	if len(c) != 0 {
		t.Error("should get an empty contour here")
	}
}

func TestNewContourMustErrorIfInputPolygonsAreNotCounterClockwiseOriented(t *testing.T) {
	clockwisePolygon := []Polygon{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}
	_, err := NewContour(clockwisePolygon)
	if err != ErrWrongOrientation {
		t.Error("should fail for wrongly oriented polygons")
	}
}

func TestNewContourReturnsInputIfInputIsASinglePolygon(t *testing.T) {
	onePolygon := []Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}}
	c, err := NewContour(onePolygon)
	if err != nil {
		t.Error("was not expecting an error here")
	}
	if len(c) != 1 {
		t.Error("was expecting exactly one polygon here")
	}
	if !EqualPolygon(c[0], onePolygon[0]) {
		t.Error("wrong polygon returned")
	}
}

func TestVerticalEdgeSortingMustSortSameAbcissaPointsLeftEdgeFirst(t *testing.T) {
	edges := []verticalEdge{}

	var sameX float64 = 42

	lastEdge := verticalEdge{sameX + 1, 2, 0}
	leftEdgeBottom := verticalEdge{sameX, 2, 0}
	leftEdgeTop := verticalEdge{sameX, 10, 5}
	rightEdge := verticalEdge{sameX, 0, 2}

	edges = append(edges, lastEdge)
	edges = append(edges, rightEdge)
	edges = append(edges, leftEdgeTop)
	edges = append(edges, leftEdgeBottom)

	sortVerticalEdges(edges)

	if !EqualEdge(edges[0], leftEdgeBottom) ||
		!EqualEdge(edges[1], leftEdgeTop) ||
		!EqualEdge(edges[2], rightEdge) ||
		!EqualEdge(edges[3], lastEdge) {
		t.Error("edge sorting is wrong")
	}
}

func TestVerticalsToHorizontals(t *testing.T) {
	testVerticals := []verticalEdge{{0.0, 7.0, 1.0}, {1.0, 1.0, 0.0}, {3.0, 0.0, 1.0},
		{5.0, 1.0, 0.0}, {6.0, 0.0, 7.0}, {2.0, 5.0, 3.0},
		{4.0, 3.0, 5.0}}

	he := verticalsToHorizontals(testVerticals)

	expected := []horizontalEdge{{1, 0, 1}, {0, 1, 3}, {1, 3, 5}, {0, 5, 6},
		{7, 6, 0}, {3, 2, 4}, {5, 4, 2}}

	if len(expected) != len(he) {
		t.Errorf("expected %d edges and got %d", len(expected), len(he))
	} else {
		for i, e := range expected {
			if !EqualEdge(e, he[i]) {
				t.Errorf("expected edge(%d) to be %s and is %s",
					i, e.String(), he[i].String())
			}
		}
	}
}

func TestFinalizeContourMustErrIfNumberOfVerticalsDifferFromNumberOfHorizontals(t *testing.T) {
	v := []verticalEdge{{0, 1, 0}, {1, 0, 1}}
	h := []horizontalEdge{{0, 0, 1}}
	_, err := finalizeContour(v, h)
	if err != errDifferentVH {
		t.Error("should have failed here")
	}
}

func checkContour(t *testing.T, contour, expected Contour) {
	if !EqualContour(contour, expected) {
		t.Error("did not get expected contour")
		fmt.Println("expected", expected)
		fmt.Println("contour", contour)
	}
}

func contourFromVerticals(t *testing.T, verticals []verticalEdge, expected Contour) {
	he := verticalsToHorizontals(verticals)
	contour, err := finalizeContour(verticals, he)
	if err != nil {
		t.Fatal("could not finalize contour")
	}
	checkContour(t, contour, expected)
}
func TestFinalizeContourIEEEExample(t *testing.T) {
	verticals := []verticalEdge{{0.0, 7.0, 1.0}, {1.0, 1.0, 0.0}, {3.0, 0.0, 1.0},
		{5.0, 1.0, 0.0}, {6.0, 0.0, 7.0}, {2.0, 5.0, 3.0},
		{4.0, 3.0, 5.0}}
	expected := Contour{
		{{0, 7}, {0, 1}, {1, 1}, {1, 0}, {3, 0}, {3, 1}, {5, 1}, {5, 0}, {6, 0}, {6, 7}, {0, 7}},
		{{2, 5}, {2, 3}, {4, 3}, {4, 5}, {2, 5}},
	}
	contourFromVerticals(t, verticals, expected)
}

func TestFinalizeContourWithOneCommonVertex(t *testing.T) {
	verticals := []verticalEdge{{0, 2, 0}, {1, 0, 2}, {1, 4, 2}, {2, 2, 4}}
	expected := Contour{{{0, 2}, {0, 0}, {1, 0}, {1, 2}, {0, 2}},
		{{1, 4}, {1, 2}, {2, 2}, {2, 4}, {1, 4}}}
	contourFromVerticals(t, verticals, expected)
}

func TestNewContourWithOneCommonVertex(t *testing.T) {
	input := []Polygon{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
		{{0, 1}, {1, 1}, {1, 2}, {0, 2}, {0, 1}},
		{{1, 2}, {2, 2}, {2, 3}, {1, 3}, {1, 2}},
		{{1, 3}, {2, 3}, {2, 4}, {1, 4}, {1, 3}},
	}
	expected := Contour{{{0, 2}, {0, 0}, {1, 0}, {1, 2}, {0, 2}},
		{{1, 4}, {1, 2}, {2, 2}, {2, 4}, {1, 4}}}

	contour, err := NewContour(input)
	if err != nil {
		t.Fatalf("could not create contour : %v", err)
	}
	checkContour(t, contour, expected)
}
