package geo

import (
	"fmt"
	"log"
)

// Contour is a set of polygons
type Contour []Polygon

// Contains returns true if the point (x,y) is inside the contour
func (c *Contour) Contains(x, y float64) bool {
	for _, p := range *c {
		inside, err := p.Contains(x, y)
		if err != nil {
			log.Fatal("got an invalid polygon") // should not happen
			// as we "control" the construction of a contour ?
			return false
		}
		if inside {
			return true
		}
	}
	return false
}

func (c *Contour) getVertices() []Vertex {
	v := []Vertex{}
	for _, p := range *c {
		v = append(v, p.getVertices()...)
	}
	return v
}

func (c *Contour) getSortedVertices() []Vertex {
	return sortVertices(c.getVertices())
}

// EqualContour returns true if the two contours have the same set of (sorted) vertices
func EqualContour(c1, c2 Contour) bool {
	v1 := c1.getSortedVertices()
	v2 := c2.getSortedVertices()

	if len(v1) != len(v2) {
		return false
	}

	if (v1 == nil) != (v2 == nil) {
		return false
	}

	for i, v := range v1 {
		if v != v2[i] {
			return false
		}
	}
	return true
}

func (c *Contour) String() string {
	s := "MULTIPOLYGON"
	for j := 0; j < len(*c); j++ {
		p := (*c)[j]
		s += "("
		for i := 0; i < len(p); i++ {
			s += fmt.Sprintf("%v %v", p[i].X, p[i].Y)
			if i < len(p)-1 {
				s += ","
			}
		}
		s += ")"
		if j < len(*c)-1 {
			s += ","
		}
	}
	s += ")"
	return s
}

// BBox returns the bounding box of this contour
func (c *Contour) BBox() BBox {
	return getVerticesBBox(c.getVertices())
}
