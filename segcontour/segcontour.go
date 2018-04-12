package segcontour

import (
	"log"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

// BBox returns the bounding box of the segmentation
func BBox(seg mapping.Segmentation) geo.BBox {
	contour := Contour(seg)
	return contour.BBox()
}

// Contour returns the contour of the segmentation
func Contour(seg mapping.Segmentation) geo.Contour {
	polygons := []geo.Polygon{}
	for _, c := range getDualSampaContours(seg) {
		for _, p := range c {
			polygons = append(polygons, p)
		}
	}
	contour, err := geo.CreateContour(polygons)
	if err != nil {
		log.Fatalf("could not get contour of segmentation: %v", err)
	}
	return contour
}

func getPadPolygons(seg mapping.Segmentation) [][]geo.Polygon {
	dualSampaPads := [][]geo.Polygon{}
	for i := 0; i < seg.NofDualSampas(); i++ {
		dualSampaPads = append(dualSampaPads, []geo.Polygon{})
		var pads []geo.Polygon
		dsID, err := seg.DualSampaID(i)
		if err != nil {
			log.Fatalf("could not get dual sampa ID: %v", err)
		}
		seg.ForEachPadInDualSampa(dsID, func(paduid int) {
			x := seg.PadPositionX(paduid)
			y := seg.PadPositionY(paduid)
			dx := seg.PadSizeX(paduid) / 2
			dy := seg.PadSizeY(paduid) / 2
			pads = append(pads, geo.Polygon{
				{X: x - dx, Y: y - dy},
				{X: x + dx, Y: y - dy},
				{X: x + dx, Y: y + dy},
				{X: x - dx, Y: y + dy},
				{X: x - dx, Y: y - dy}})
		})
		dualSampaPads[len(dualSampaPads)-1] = append(dualSampaPads[len(dualSampaPads)-1], pads...)
	}
	return dualSampaPads
}

func getDualSampaContours(seg mapping.Segmentation) []geo.Contour {
	contours := []geo.Contour{}
	padPolygons := getPadPolygons(seg)
	for _, p := range padPolygons {
		c, err := geo.CreateContour(p)
		if err != nil {
			log.Fatalf("could not create contour: %v", err)
		}
		contours = append(contours, c)
	}
	return contours
}
