package geo

import "testing"

func TestContoursAreEqualAsLongAsTheyContainTheSameSetOfVertices(t *testing.T) {
	aCollectionWithOnePolygon := Contour{{{0, 2}, {0, 0}, {2, 0}, {2, 4}, {1, 4}, {1, 2}, {0, 2}}}

	anotherCollectionWithTwoPolygonsButSameVertices := Contour{
		{{2, 4}, {2, 0}}, {{1, 4}, {1, 2}, {0, 2}, {0, 0}}}

	if !EqualContour(aCollectionWithOnePolygon, anotherCollectionWithTwoPolygonsButSameVertices) {
		t.Error("contours should be equal")
	}
}
