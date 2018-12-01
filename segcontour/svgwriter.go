package segcontour

import (
	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
)

// ShowFlags is a option struct to select what to show / hide in SVG outputs
type ShowFlags struct {
	des         bool
	dualsampas  bool
	pads        bool
	padchannels bool
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
	if show.des {
		deContour := Contour(cseg)
		svgDetectionElements(w, &deContour)
	}
	if show.pads {
		dualSampaPads := getAllDualSampaPadPolygons(cseg)
		svgDualSampaPads(w, &dualSampaPads)
	}
}
