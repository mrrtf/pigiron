package main

import (
	"fmt"
	"os"

	"github.com/aphecetche/pigiron/geo"
	"github.com/aphecetche/pigiron/mapping"

	// must include the specific implementation package of the mapping
	_ "github.com/aphecetche/pigiron/mapping/impl4"
)

func main() {

	type padSize struct {
		x, y float64
	}
	padsizes := make(map[padSize]int)

	mapping.ForOneDetectionElementOfEachSegmentationType(func(deid int) {
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
