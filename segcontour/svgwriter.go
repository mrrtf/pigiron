package segcontour

import (
	"fmt"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

// ShowFlags is a option struct to select what to show / hide in SVG outputs
type ShowFlags struct {
	DEs         bool
	DualSampas  bool
	Pads        bool
	PadChannels bool
}

func svgDualSampaPads(w *geo.SVGWriter, dualSampaPads *[][]geo.Polygon) {
	w.GroupStart("pads")
	defer w.GroupEnd()
	for _, dsp := range *dualSampaPads {
		for _, p := range dsp {
			w.Polygon(&p)
		}
	}
}

func svgDetectionElements(w *geo.SVGWriter, de *geo.Contour) {
	w.GroupStart("detectionelements")
	defer w.GroupEnd()
	w.Contour(de)
}

// SVGSegmentation creates a SVG representation of segmentation
func SVGSegmentation(cseg mapping.CathodeSegmentation, w *geo.SVGWriter, show ShowFlags) {
	if show.Pads {
		dualSampaPads := getAllDualSampaPadPolygons(cseg)
		svgDualSampaPads(w, &dualSampaPads)
	}
	if show.DEs {
		deContour := Contour(cseg)
		svgDetectionElements(w, &deContour)
	}
	if show.DualSampas {
		w.GroupStart("dualsampas")
		for i := 0; i < cseg.NofDualSampas(); i++ {
			dsid, err := cseg.DualSampaID(i)
			if err != nil {
				panic(err)
			}
			dsContour := GetDualSampaContour(cseg, dsid)
			for _, c := range dsContour {
				w.PolygonWithClass(&c, fmt.Sprintf("polds DS%d", dsid))
			}
		}
		w.GroupEnd()
	}
}
