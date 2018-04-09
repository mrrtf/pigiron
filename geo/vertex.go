package geo

import (
	"fmt"
	"log"
	"math"
)

// Vertex is a simple (x,y) float64 pair
type Vertex struct {
	x, y float64
}

func isVerticalSegment(a, b Vertex) bool {
	return EqualFloat(a.x, b.x)
}

func isHorizontalSegment(a, b Vertex) bool {
	return EqualFloat(a.y, b.y)
}

func (v Vertex) String() string {
	return fmt.Sprintf("( %v %v )", v.x, v.y)
}

func sub(a, b Vertex) Vertex {
	return Vertex{a.x - b.x, a.y - b.y}
}

// EqualVertex checks if two vertices are equal
// For the precision of the comparison see EqualFloat function.
func EqualVertex(a, b Vertex) bool {
	return EqualFloat(a.x, b.x) &&
		EqualFloat(a.y, b.y)
}

func dot(a, b Vertex) float64 {
	return a.x*b.x + a.y*b.y
}

func squaredDistance(a, b Vertex) float64 {
	aminusb := sub(a, b)
	return dot(aminusb, aminusb)
}

//SquaredDistanceOfPointToSegment computes the square of the distance
// between point p and segment (p0,p1)
func SquaredDistanceOfPointToSegment(p, p0, p1 Vertex) float64 {
	v := sub(p1, p0)
	w := sub(p, p0)
	c1 := dot(w, v)

	if c1 <= 0 {
		return squaredDistance(p, p0)
	}

	c2 := dot(v, v)
	if c2 <= c1 {
		return squaredDistance(p, p1)
	}

	b := c1 / c2
	pbase := Vertex{p0.x + b*v.x, p0.y + b*v.y}
	return squaredDistance(p, pbase)
}

func getVerticesBBox(vertices []Vertex) BBox {
	xmin := math.MaxFloat64
	xmax := -xmin
	ymin := xmin
	ymax := -ymin

	for _, v := range vertices {
		xmin = math.Min(xmin, v.x)
		ymin = math.Min(ymin, v.y)
		xmax = math.Max(xmax, v.x)
		ymax = math.Max(ymax, v.y)
	}
	bbox, err := NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		log.Fatal("got a very unexpected invalid bbox here")
	}
	return bbox
}
