package geo

import (
	"fmt"
	"log"
	"math"
)

// Vertex is a simple (x,y) float64 pair
type Vertex struct {
	X, Y float64
}

func isVerticalSegment(a, b Vertex) bool {
	return EqualFloat(a.X, b.X)
}

func isHorizontalSegment(a, b Vertex) bool {
	return EqualFloat(a.Y, b.Y)
}

func (v Vertex) String() string {
	return fmt.Sprintf("( %v %v )", v.X, v.Y)
}

func sub(a, b Vertex) Vertex {
	return Vertex{a.X - b.X, a.Y - b.Y}
}

// EqualVertex checks if two vertices are equal
// For the precision of the comparison see EqualFloat function.
func EqualVertex(a, b Vertex) bool {
	return EqualFloat(a.X, b.X) &&
		EqualFloat(a.Y, b.Y)
}

func dot(a, b Vertex) float64 {
	return a.X*b.X + a.Y*b.Y
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
	pbase := Vertex{p0.X + b*v.X, p0.Y + b*v.Y}
	return squaredDistance(p, pbase)
}

func getVerticesBBox(vertices []Vertex) BBox {
	xmin := math.MaxFloat64
	xmax := -xmin
	ymin := xmin
	ymax := -ymin

	for _, v := range vertices {
		xmin = math.Min(xmin, v.X)
		ymin = math.Min(ymin, v.Y)
		xmax = math.Max(xmax, v.X)
		ymax = math.Max(ymax, v.Y)
	}
	bbox, err := NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		log.Fatal("got a very unexpected invalid bbox here")
	}
	return bbox
}
