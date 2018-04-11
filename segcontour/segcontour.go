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
	var polygons []geo.Polygon
	for _, c := range getAllDualSampaContours(seg) {
		for _, p := range c {
			polygons = append(polygons, p)
		}
	}
	contour, err := geo.NewContour(polygons)
	if err != nil {
		log.Fatalf("could not get contour of segmentation: %v", err)
	}
	return contour
}

func getDualSampaPadPolygons(seg mapping.Segmentation, dualSampaID int) []geo.Polygon {
	var pads []geo.Polygon
	seg.ForEachPadInDualSampa(dualSampaID, func(paduid int) {
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
	return pads
}

// GetDualSampaContour returns the contour of one FEC.
func GetDualSampaContour(seg mapping.Segmentation, dualSampaID int) geo.Contour {
	pads := getDualSampaPadPolygons(seg, dualSampaID)
	c, err := geo.NewContour(pads)
	if err != nil {
		log.Fatalf("could not create contour : %v", err)
	}
	return c
}

func getAllDualSampaPadPolygons(seg mapping.Segmentation) [][]geo.Polygon {
	dualSampaPads := [][]geo.Polygon{}
	for i := 0; i < seg.NofDualSampas(); i++ {
		dualSampaPads = append(dualSampaPads, []geo.Polygon{})
		dsID, err := seg.DualSampaID(i)
		if err != nil {
			log.Fatalf("could not get dual sampa ID: %v", err)
		}
		dualSampaPads[len(dualSampaPads)-1] = append(dualSampaPads[len(dualSampaPads)-1], getDualSampaPadPolygons(seg, dsID)...)
	}
	return dualSampaPads
}

func getAllDualSampaContours(seg mapping.Segmentation) []geo.Contour {
	var contours = []geo.Contour{}
	for i := 0; i < seg.NofDualSampas(); i++ {
		dsID, err := seg.DualSampaID(i)
		var c geo.Contour
		if err == nil {
			c = GetDualSampaContour(seg, dsID)
		}
		contours = append(contours, c)
	}
	return contours
}
