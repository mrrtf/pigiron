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
	DEID     int      `json:"deid"`
	DSID     int      `json:"dsid"`
	DSCH     int      `json:"dsch"`
	Vertices []Vertex `json:"vertices"`
}

type PadAlt struct {
	DEID     int      `json:"deid"`
	PADID    int      `json:"padid"`
	KEY      string   `json:"id"`
	Vertices []Vertex `json:"vertices"`
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

type FlipDirection int

const (
	FlipNone FlipDirection = iota + 1
	FlipX
	FlipY
	FlipXY
)

func flipVertex(p Vertex, d FlipDirection, xcenter, ycenter float64) Vertex {
	if d == FlipX || d == FlipXY {
		p.Y = 2*ycenter - p.Y
	}
	if d == FlipY || d == FlipXY {
		p.X = 2*xcenter - p.X
	}
	return p
}

var de2rot = map[int]FlipDirection{}

func initDetectionElementRotations() {
	for i := 1; i <= 4; i++ {
		de2rot[i*100] = FlipX
		de2rot[i*100+1] = FlipXY
		de2rot[i*100+2] = FlipY
		de2rot[i*100+3] = FlipNone
	}
	for i := 0; i < 2; i++ {
		de2rot[500+i*100] = FlipXY
		de2rot[501+i*100] = FlipXY
		de2rot[502+i*100] = FlipNone
		de2rot[503+i*100] = FlipX
		de2rot[504+i*100] = FlipNone
		de2rot[505+i*100] = FlipY
		de2rot[506+i*100] = FlipXY
		de2rot[507+i*100] = FlipY
		de2rot[508+i*100] = FlipX
		de2rot[509+i*100] = FlipX
		de2rot[510+i*100] = FlipNone
		de2rot[511+i*100] = FlipY
		de2rot[512+i*100] = FlipXY
		de2rot[513+i*100] = FlipY
		de2rot[514+i*100] = FlipNone
		de2rot[515+i*100] = FlipX
		de2rot[516+i*100] = FlipNone
		de2rot[517+i*100] = FlipY
	}
	for i := 0; i < 4; i++ {
		de2rot[700+i*100] = FlipNone
		de2rot[701+i*100] = FlipXY
		de2rot[702+i*100] = FlipNone
		de2rot[703+i*100] = FlipX
		de2rot[704+i*100] = FlipNone
		de2rot[705+i*100] = FlipX
		de2rot[706+i*100] = FlipNone
		de2rot[707+i*100] = FlipY
		de2rot[708+i*100] = FlipXY
		de2rot[709+i*100] = FlipY
		de2rot[710+i*100] = FlipXY
		de2rot[711+i*100] = FlipY
		de2rot[712+i*100] = FlipX
		de2rot[713+i*100] = FlipY
		de2rot[714+i*100] = FlipNone
		de2rot[715+i*100] = FlipY
		de2rot[716+i*100] = FlipXY
		de2rot[717+i*100] = FlipY
		de2rot[718+i*100] = FlipXY
		de2rot[719+i*100] = FlipY
		de2rot[720+i*100] = FlipNone
		de2rot[721+i*100] = FlipX
		de2rot[722+i*100] = FlipNone
		de2rot[723+i*100] = FlipX
		de2rot[724+i*100] = FlipNone
		de2rot[725+i*100] = FlipY
	}
}

func init() {
	initDetectionElementRotations()
}

func flipVertices(vertices []Vertex, deid int, xcenter float64, ycenter float64) []Vertex {
	var flipped []Vertex
	direction, _ := de2rot[deid]
	for _, v := range vertices {
		flipped = append(flipped, flipVertex(v, direction, xcenter, ycenter))
	}
	return flipped
}

func jsonDEGeo(w io.Writer, cseg mapping.CathodeSegmentation, bending bool) {

	bbox := mapping.ComputeBBox(cseg)

	var vertices []Vertex
	contour := segcontour.Contour(cseg)
	deid := int(cseg.DetElemID())

	for _, c := range contour {
		for _, v := range c {
			vertices = append(vertices, Vertex{X: v.X, Y: v.Y})
		}
	}

	degeo := DEGeo{
		ID:      deid,
		Bending: cseg.IsBending(),
		X:       bbox.Xcenter(),
		Y:       bbox.Ycenter(),
		SX:      bbox.Width(),
		SY:      bbox.Height()}

	degeo.Vertices = flipVertices(vertices, deid, bbox.Xcenter(), bbox.Ycenter())
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

	bbox := mapping.ComputeBBox(cseg)

	deid := int(cseg.DetElemID())

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

		ds.Vertices = flipVertices(ds.Vertices, deid, bbox.Xcenter(), bbox.Ycenter())

		dualSampas = append(dualSampas, ds)
	}

	b, err := json.Marshal(dualSampas)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprintf(w, string(b))
}

func jsonPadList(w io.Writer, padlist []PadRef, keepOrder bool) {

	if !keepOrder {
		jsonPadListRegular(w, padlist)
	} else {
		jsonPadListOrdered(w, padlist)
	}
}

type Vertices struct {
	Vertices []Vertex `json:"vertices"`
}

func createVertexer(segcache mapping.SegCache) func(pad PadRef) (string, []Vertex) {
	return func(pad PadRef) (string, []Vertex) {
		paduid := mapping.PadUID(pad.PadId)
		seg := segcache.Segmentation(mapping.DEID(pad.DeId))
		cseg := seg.Bending()
		if !seg.IsBendingPad(paduid) {
			cseg = seg.NonBending()
		}
		bbox := mapping.ComputeBBox(cseg)
		x := seg.PadPositionX(paduid)
		y := seg.PadPositionY(paduid)
		dx := seg.PadSizeX(paduid) / 2
		dy := seg.PadSizeY(paduid) / 2
		var vertices []Vertex
		vertices = append(vertices, Vertex{X: x - dx, Y: y - dy})
		vertices = append(vertices, Vertex{X: x + dx, Y: y - dy})
		vertices = append(vertices, Vertex{X: x + dx, Y: y + dy})
		vertices = append(vertices, Vertex{X: x - dx, Y: y + dy})
		vertices = append(vertices, Vertex{X: x - dx, Y: y - dy})
		key := fmt.Sprintf("%v-%v-%v", pad.DeId,
			int(seg.PadDualSampaID(paduid)),
			int(seg.PadDualSampaChannel(paduid)))
		vertices = flipVertices(vertices, pad.DeId, bbox.Xcenter(), bbox.Ycenter())
		return key, vertices
	}
}

func jsonPadListRegular(w io.Writer, padlist []PadRef) {
	pads := make(map[string]Vertices)

	segcache := mapping.SegCache{}

	vertexer := createVertexer(segcache)

	for _, pad := range padlist {
		key, vertices := vertexer(pad)
		pads[key] = Vertices{vertices}
	}

	b, err := json.Marshal(pads)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprintf(w, string(b))
}

// type PadAlt struct {
// 	DEID     int      `json:"deid"`
// 	PADID    int      `json:"padid"`
// 	KEY      string   `json:"key"`
// 	Vertices []Vertex `json:"vertices"`
// }

func jsonPadListOrdered(w io.Writer, padlist []PadRef) {
	var pads []PadAlt

	segcache := mapping.SegCache{}

	vertexer := createVertexer(segcache)

	for _, pad := range padlist {
		key, vertices := vertexer(pad)
		pads = append(pads, PadAlt{DEID: pad.DeId, PADID: pad.PadId, KEY: key, Vertices: vertices})
	}

	b, err := json.Marshal(pads)

	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprintf(w, string(b))
}
