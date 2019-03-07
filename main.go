package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"
	"github.com/aphecetche/pigiron/segcontour"

	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

type testarea struct {
	xcenter, ycenter float64
	width, height    float64
}

func main() {

	plotSizes()

	plotTestAreas(100, []testarea{
		{52.50, 20.0, 5.0, 1.0}})
	// {24.0, 24.0, 0.45, 0.45},
	// {5.7, 17.5, 9.3, 4.3},
	// {15.5, 10.0, 5.5, 10.0},
	// {85.5, 40.0, 6.0, 20.0}})

	plotTestAreas(500, []testarea{
		{51, 0.0, 12.0, 35.0}})

	plotTestAreas(701, []testarea{
		{99, -7.0, 12.0, 35.0}})
}

func plotTestAreas(deid int, areas []testarea) {

	cssStyle := `
.pads {
  fill: #EEEEEE;
  stroke-width: 0.025px;
  stroke: #AAAAAA;
}
.dualsampas {
  fill:none;
  stroke-width: 0.025px;
  stroke: #333333;
}
.detectionelements {
  fill:none;
  stroke-width:0.025px;
  stroke: #000000;
}
.center {
  fill:red;
  stroke-width:0.025px;
  stroke: black;
  opacity: 0.5;
}
.area {
  fill:blue;
  opacity: 0.5;
}
.padsinarea {
stroke: blue;
stroke-width: 0.05px;
fill: yellow;
opacity: 0.5;
}
.clippedArea {
fill:pink;
opacity:0.5;
stroke:red;
stroke-width: 0.075px;
}
`

	for _, bending := range []bool{true, false} {
		cseg := mapping.NewCathodeSegmentation(mapping.DEID(deid), bending)
		svg := geo.NewSVGWriter(1024)
		svg.Style(cssStyle)
		showflags := segcontour.ShowFlags{DEs: true, Pads: true, PadChannels: true}
		segcontour.SVGSegmentation(cseg, svg, showflags)

		for _, a := range areas {

			plotArea(cseg, svg, a.xcenter-a.width/2.0, a.ycenter-a.height/2.0, a.xcenter+a.width/2.0, a.ycenter+a.height/2.0)
		}

		f, _ := os.Create(fmt.Sprintf("segmentation-%d-%s.html", deid, mapping.PlaneAbbreviation(bending)))

		svg.MoveToOrigin()
		svg.WriteHTML(f)
	}
}

func plotSizes() {

	type padSize struct {
		x, y float64
	}
	padsizes := make(map[padSize]int)

	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid mapping.DEID) {
		for _, isBending := range []bool{true, false} {
			seg := mapping.NewCathodeSegmentation(deid, isBending)
			seg.ForEachPad(func(padcid mapping.PadCID) {
				ps := &padSize{seg.PadSizeX(padcid), seg.PadSizeY(padcid)}
				padsizes[*ps]++
			})
		}
	})

	if len(padsizes) != 18 {
		fmt.Errorf("wanted 18 padsizes - got %d", len(padsizes))
	}

	svg := geo.NewSVGWriter(1024)

	svg.Style(`
rect {
stroke: red;
stroke-width: 0.02;
fill: none;
}`)

	for ps := range padsizes {
		svg.Rect(1.0, 1.0, ps.x, ps.y)
	}

	f, _ := os.Create("padsizes.html")
	svg.WriteHTML(f)
}

func plotArea(cseg mapping.CathodeSegmentation, svg *geo.SVGWriter, xmin, ymin, xmax, ymax float64) {
	x := (xmin + xmax) / 2.0
	y := (ymin + ymax) / 2.0
	svg.CircleWithClass(x, y, 0.05, "center")
	svg.RectWithClass(xmin, ymin, (xmax - xmin), (ymax - ymin), "area")

	svg.GroupStart("padsinarea")
	cseg.ForEachPadInArea(xmin, ymin, xmax, ymax, func(padcid mapping.PadCID) {
		var xpadmin, ypadmin, xpadmax, ypadmax float64
		mapping.ComputeCathodePadBBox(cseg, padcid, &xpadmin, &ypadmin, &xpadmax, &ypadmax)
		svg.RectWithClass(xpadmin, ypadmin, (xpadmax - xpadmin), (ypadmax - ypadmin), "padsinarea")
	})
	svg.GroupEnd()

	segContour := segcontour.Contour(cseg)

	bbox, err := geo.NewBBox(xmin, ymin, xmax, ymax)
	if err != nil {
		log.Fatalf("Could not create bbox : %s\n,", err.Error())
	}

	pol, err := geo.ClipPolygon(segContour[0], bbox)
	if err != nil {
		panic(err)
	}
	svg.PolygonWithClass(&pol, "clippedArea")

}
