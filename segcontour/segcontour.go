package segcontour

import (
	"log"

	"github.com/mrrtf/pigiron/geo"
	"github.com/mrrtf/pigiron/mapping"
)

// BBox returns the bounding box of the segmentation
func BBox(cseg mapping.CathodeSegmentation) geo.BBox {
	contour := Contour(cseg)
	return contour.BBox()
}

// Contour returns the contour of the segmentation
func Contour(cseg mapping.CathodeSegmentation) geo.Contour {
	var polygons []geo.Polygon
	for _, c := range getAllDualSampaContours(cseg) {
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

func GetDualSampaPadPolygons(cseg mapping.CathodeSegmentation, dsid mapping.DualSampaID) []geo.Polygon {
	var pads []geo.Polygon
	cseg.ForEachPadInDualSampa(dsid, func(padcid mapping.PadCID) {
		x := cseg.PadPositionX(padcid)
		y := cseg.PadPositionY(padcid)
		dx := cseg.PadSizeX(padcid) / 2
		dy := cseg.PadSizeY(padcid) / 2
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
func GetDualSampaContour(cseg mapping.CathodeSegmentation, dsid mapping.DualSampaID) geo.Contour {
	pads := GetDualSampaPadPolygons(cseg, dsid)
	c, err := geo.NewContour(pads)
	if err != nil {
		log.Fatalf("could not create contour : %v", err)
	}
	return c
}

func getAllDualSampaPadPolygons(cseg mapping.CathodeSegmentation) [][]geo.Polygon {
	dualSampaPads := [][]geo.Polygon{}
	for i := 0; i < cseg.NofDualSampas(); i++ {
		dualSampaPads = append(dualSampaPads, []geo.Polygon{})
		dsID, err := cseg.DualSampaID(i)
		if err != nil {
			log.Fatalf("could not get dual sampa ID: %v", err)
		}
		dualSampaPads[len(dualSampaPads)-1] = append(dualSampaPads[len(dualSampaPads)-1], GetDualSampaPadPolygons(cseg, dsID)...)
	}
	return dualSampaPads
}

func getAllDualSampaContours(cseg mapping.CathodeSegmentation) []geo.Contour {
	var contours = []geo.Contour{}
	for i := 0; i < cseg.NofDualSampas(); i++ {
		dsID, err := cseg.DualSampaID(i)
		var c geo.Contour
		if err == nil {
			c = GetDualSampaContour(cseg, dsID)
		}
		contours = append(contours, c)
	}
	return contours
}
