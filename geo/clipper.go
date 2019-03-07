package geo

import (
	"fmt"
	"log"
)

// ClipPolygon will returns the polygon resulting from
// the clipping of pol with window.
//
// Algorithm from "Reetrant Polygon Clipping",
// Ivan E. Sutherland and Gary W. Hodgman
// Communications of the ACM, January 1974, Volume 17, Number 1, p32-42
//
func ClipPolygon(pol Polygon, w BBox) (Polygon, error) {

	if !pol.isCounterClockwiseOriented() {
		return Polygon{}, ErrWrongOrientation
	}

	// clipWindow, counter-clockwise oriented
	clipWindow := Polygon{
		Vertex{w.Xmax(), w.Ymax()},
		Vertex{w.Xmin(), w.Ymax()},
		Vertex{w.Xmin(), w.Ymin()},
		Vertex{w.Xmax(), w.Ymin()}}

	clipWindow = ClosePolygon(clipWindow)

	var clipEdges []Segment

	for i := 0; i < len(clipWindow)-1; i++ {
		clipEdges = append(clipEdges, Segment{clipWindow[i], clipWindow[i+1]})
	}

	p := ClosePolygon(pol)
	for _, ce := range clipEdges {
		p = clipWithOneEdge(p, ce)
	}
	return p, nil
}

func polAsString(pol Polygon, clipEdge Segment) string {
	s := "("
	for _, v := range pol {
		s += fmt.Sprintf("%4.0f %4.0f", v.X, v.Y)
		if isInside(v, clipEdge) {
			s += " [],"
		} else {
			s += " x,"
		}
	}
	s += ")"
	return s
}

func clipWithOneEdge(pol Polygon, clipEdge Segment) Polygon {
	// fmt.Printf("\n>>> Clipping polygon %v\n>>> with edge %v\n", polAsString(pol, clipEdge), clipEdge)
	var output []Vertex
	s := pol[0]
	if isInside(s, clipEdge) {
		output = append(output, s)
	}
	for i := 1; i < len(pol); i++ {
		e := pol[i]
		if isInside(e, clipEdge) {
			if !isInside(s, clipEdge) {
				i, ok := IntersectSegmentLine(Segment{s, e}, clipEdge)
				if !ok {
					log.Fatalf("got a parallel for case 1 ?")
				}
				output = append(output, i)
			}
			output = append(output, e)
		} else if isInside(s, clipEdge) {
			i, ok := IntersectSegmentLine(Segment{s, e}, clipEdge)
			if !ok {
				log.Fatalf("got a parallel for case 2 se=%v clipEdge=%v\n", Segment{s, e}, clipEdge)
			}
			output = append(output, i)
		}
		s = e
	}
	return ClosePolygon(output)
}

func isInside(v Vertex, s Segment) bool {
	p := (s.P1.X - s.P0.X) * (v.Y - s.P0.Y)
	p -= (s.P1.Y - s.P0.Y) * (v.X - s.P0.X)
	return p > 0
}
