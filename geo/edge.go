package geo

import (
	"fmt"
	"math"
)

type manhattanEdge interface {
	begin() Vertex
	end() Vertex
}

type verticalEdge struct {
	x, y1, y2 float64
}

type horizontalEdge struct {
	y, x1, x2 float64
}

func NewEdge(x1, y1, x2, y2 float64) (manhattanEdge, error) {
	v1 := Vertex{x1, y1}
	v2 := Vertex{x2, y1}
	if isVerticalSegment(v1, v2) {
		return verticalEdge{x1, y1, y2}, nil
	}
	if isHorizontalSegment(v1, v2) {
		return horizontalEdge{y1, x1, x2}, nil
	}
	return nil, fmt.Errorf("coordinates incompatible with a manhattan edge")
}

func (v verticalEdge) begin() Vertex {
	return Vertex{v.x, v.y1}
}

func (v verticalEdge) end() Vertex {
	return Vertex{v.x, v.y2}
}

func (h horizontalEdge) begin() Vertex {
	return Vertex{h.x1, h.y}
}

func (h horizontalEdge) end() Vertex {
	return Vertex{h.x2, h.y}
}

func isLeftEdge(v verticalEdge) bool {
	return v.begin().Y > v.end().Y
}

func isTopToBottom(v verticalEdge) bool {
	return isLeftEdge(v)
}

func isRightEdge(v verticalEdge) bool {
	return !isLeftEdge(v)
}

func isBottomToTop(v verticalEdge) bool {
	return !isTopToBottom(v)
}

func isLeftToRight(h horizontalEdge) bool {
	return h.begin().X < h.end().X
}

func isRightToLeft(h horizontalEdge) bool {
	return !isLeftToRight(h)
}

func top(v verticalEdge) float64 {
	return math.Max(v.begin().Y, v.end().Y)
}

func bottom(v verticalEdge) float64 {
	return math.Min(v.begin().Y, v.end().Y)
}

func left(h horizontalEdge) float64 {
	return math.Min(h.begin().X, h.end().X)
}

func right(h horizontalEdge) float64 {
	return math.Max(h.begin().X, h.end().X)
}

func (h horizontalEdge) String() string {
	s := fmt.Sprintf("[%v,%v]", h.begin(), h.end())
	if isLeftToRight(h) {
		s += "->-"
	} else {
		s += "-<-"
	}
	return s
}

func (v verticalEdge) String() string {
	s := fmt.Sprintf("[%v,%v]", v.begin(), v.end())
	if isTopToBottom(v) {
		s += "v"
	} else {
		s += "^"
	}
	return s

}

// EqualEdge returns true if edges a and b are equal
func EqualEdge(a, b manhattanEdge) bool {
	return EqualVertex(a.begin(), b.begin()) &&
		EqualVertex(a.end(), b.end())
}
