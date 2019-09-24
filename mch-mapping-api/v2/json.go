package v2

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/mrrtf/pigiron/mapping"
	"github.com/mrrtf/pigiron/segcontour"
)

type Vertex struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Pad struct {
	DSID int     `json:"DSID"`
	DSCH int     `json:"DSCH"`
	X    float64 `json:"X"`
	Y    float64 `json:"Y"`
	SX   float64 `json:"SX"`
	SY   float64 `json:"SY"`
}

type DualSampaPads struct {
	ID   int `json:"id"`
	Pads []Pad
}

type DualSampa struct {
	ID       int      `json:"id"`
	Vertices []Vertex `json:"vertices"`
}

type DEGeo struct {
	ID       int      `json:"id"`
	Bending  bool     `json:"bending"`
	X        float64  `json:"x"`
	Y        float64  `json:"y"`
	SX       float64  `json:"sx"`
	SY       float64  `json:"sy"`
	Vertices []Vertex `json:"vertices"`
}

func jsonDEGeo(w io.Writer, cseg mapping.CathodeSegmentation, bending bool) {

	bbox := mapping.ComputeBBox(cseg)

	var vertices []Vertex
	contour := segcontour.Contour(cseg)
	for _, c := range contour {
		for _, v := range c {
			vertices = append(vertices, Vertex{X: v.X, Y: v.Y})
		}
	}

	degeo := DEGeo{
		ID:      int(cseg.DetElemID()),
		Bending: cseg.IsBending(),
		X:       bbox.Xcenter(),
		Y:       bbox.Ycenter(),
		SX:      bbox.Width(),
		SY:      bbox.Height()}

	degeo.Vertices = vertices
	b, err := json.Marshal(degeo)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprintf(w, string(b))
}

func jsonDualSampaPads(w io.Writer, cseg mapping.CathodeSegmentation, dsid int) {
	var dualSampas []DualSampa
	n := cseg.NofDualSampas()

	for i := 0; i < n; i++ {
		dsid, err := cseg.DualSampaID(i)
		if err != nil {
			panic(err)
		}

		ds := DualSampa{ID: int(dsid)}

		padContours := segcontour.GetDualSampaPadPolygons(cseg, mapping.DualSampaID(dsid))

		for _, c := range padContours {
			for _, v := range c {
				ds.Vertices = append(ds.Vertices, Vertex{X: v.X, Y: v.Y})
			}
		}

		dualSampas = append(dualSampas, ds)
	}

	b, err := json.Marshal(dualSampas)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprintf(w, string(b))
}

func jsonDualSampas(w io.Writer, cseg mapping.CathodeSegmentation, bending bool) {

	var dualSampas []DualSampa
	n := cseg.NofDualSampas()

	for i := 0; i < n; i++ {
		dsid, err := cseg.DualSampaID(i)
		if err != nil {
			panic(err)
		}

		ds := DualSampa{ID: int(dsid)}

		dsContour := segcontour.GetDualSampaContour(cseg, dsid)
		for _, c := range dsContour {
			for _, v := range c {
				ds.Vertices = append(ds.Vertices, Vertex{X: v.X, Y: v.Y})
			}
		}

		dualSampas = append(dualSampas, ds)
	}

	b, err := json.Marshal(dualSampas)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprintf(w, string(b))
}
